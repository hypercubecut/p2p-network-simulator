package bitcoin

import (
	"fmt"
	"math"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/bitcoin/servicecode"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/bytedance/ns-x/v2/base"
)

var (
	baseDelay = time.Duration(2)
	count     = 0
)

func (n *Node) handleMineNewBlock(request *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	mclock.Lock()
	defer mclock.Unlock()

	if n.serviceCode != servicecode.MinerNode {
		return nil
	}

	n.logger.Debug(fmt.Sprintf("%s start mine a new block", n.name), zap.Int("chainLen", len(n.chain)))

	prevLastBlock := MasterBlockchain[len(MasterBlockchain)-1]

	reqDTO, ok := request.Payload.(*WriteBlockRequest)
	if !ok {
		n.logger.Error("handleWriteBlock failed unmarshal payload")
		return n.handleErrResp(msgtype.MineNewBlockResp, ErrUnknownPayload, request)
	}

	newBlock, err := GenerateBlock(prevLastBlock, reqDTO.BPM)
	if err != nil {
		n.logger.Error(fmt.Sprintf("%s failed to generateBlock", n.name), zap.Error(err))
		return n.handleErrResp(msgtype.MineNewBlockResp, err, request)
	}

	MasterBlockchain = append(MasterBlockchain, newBlock)

	n.newMinedBlock = newBlock
	n.logger.Info(fmt.Sprintf("%s finish mine a new block", n.name), zap.Int("newBlock", newBlock.Index))

	return n.handleStandardBlockRelay(newBlock, nodes, now)
}

func (n *Node) handleMineSameBlock(request *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	mclock.Lock()
	defer mclock.Unlock()

	if n.serviceCode != servicecode.MinerNode {
		return nil
	}

	n.logger.Debug(fmt.Sprintf("%s start mine a new block", n.name), zap.Int("chainLen", len(n.chain)))

	prevLastBlock := MasterBlockchain[len(MasterBlockchain)-1]

	reqDTO, ok := request.Payload.(*WriteBlockRequest)
	if !ok {
		n.logger.Error("handleWriteBlock failed unmarshal payload")
		return n.handleErrResp(msgtype.MineNewBlockResp, ErrUnknownPayload, request)
	}

	newBlock, err := GenerateBlock(prevLastBlock, reqDTO.BPM)
	if err != nil {
		n.logger.Error(fmt.Sprintf("%s failed to generateBlock", n.name), zap.Error(err))
		return n.handleErrResp(msgtype.MineNewBlockResp, err, request)
	}

	//MasterBlockchain = append(MasterBlockchain, newBlock)

	n.newMinedBlock = newBlock
	//n.logger.Info(fmt.Sprintf("%s finish mine a new block", n.name), zap.Int("newBlock", newBlock.Index))

	return n.handleStandardBlockRelay(newBlock, nodes, now)
}

func getIDFromName(name string) int64 {
	idStr := strings.Split(name, "-")[1]
	id, _ := strconv.ParseInt(idStr, 0, 64)

	return id
}

func (n *Node) handleStandardBlockRelay(newBlock *Block, nodes map[string]base.Node, now time.Time) []base.Event {
	return n.sendPacketToPeerNodes(&Packet{
		MessageType: msgtype.InventoryMessage,
		Payload:     &InventoryMessage{Inventory: newBlock.Index},
		Source:      n,
	}, nodes, now)
}

func (n *Node) sendPacketToPeerNodes(packet *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	var events []base.Event

	myID := getIDFromName(n.name)

	for peer, _ := range n.availablePeers {
		node := nodes[peer]
		dest, ok := node.(*Node)
		if ok {
			if dest.name == n.name ||
				dest.serviceCode != servicecode.FullNode {
				continue
			}

			destID := getIDFromName(dest.name)

			//if math.Abs(float64(myID-destID)) > numberOfPeers {
			//	continue
			//}

			event := n.Send(&Packet{
				MessageType: packet.MessageType,
				Payload:     packet.Payload,
				Source:      packet.Source,
				Destination: dest,
				SizeInBytes: packet.SizeInBytes,
			}, now.Add(calculatePeerDelay(myID, destID)))

			//fmt.Println(n.name, "send to peer", dest.name, now)

			events = append(events, event)
		}
	}

	return events
}

func calculatePeerDelay(myID, destID int64) time.Duration {
	return time.Duration(math.Abs(float64(myID-destID))) * time.Millisecond
}

func (n *Node) handleInventoryMessage(packet *Packet, now time.Time) []base.Event {

	switch concrete := packet.Payload.(type) {
	case *InventoryMessage:
		// full node request for the new block if valid
		if n.serviceCode == servicecode.FullNode && concrete.Inventory == n.inventory+1 {
			var events []base.Event

			events = append(events, n.Send(&Packet{
				MessageType: msgtype.GetNewBlockMessageType,
				Payload: &GetBlockDataReq{
					Index: concrete.Inventory,
				},
				Source:      n,
				Destination: packet.Source,
			}, now))

			//fmt.Println(n.name, "send getBlock message to", packet.Source.name)

			return events
		}

		if concrete.Inventory > n.inventory+1 {
			event := n.Send(&Packet{
				MessageType: msgtype.InvalidInventoryMessage,
				Payload:     nil,
				Source:      n,
				Destination: packet.Source,
			}, now)

			return base.Aggregate(event)
		}

		return nil

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) handleInvalidInventory(nodes map[string]base.Node, now time.Time) []base.Event {
	if n.serviceCode == servicecode.MinerNode {
		//fmt.Println(n.name, "retry publish block")
		return n.handleStandardBlockRelay(n.newMinedBlock, nodes, now.Add(5*time.Second))
	}

	return nil
}

const blockSize = 2000000

func (n *Node) getNewBlockDataHandler(packet *Packet, now time.Time) []base.Event {
	var event base.Event

	switch concrete := packet.Payload.(type) {
	case *GetBlockDataReq:
		if n.serviceCode == servicecode.FullNode {
			if concrete.Index >= len(n.chain) {
				fmt.Println(n.name, concrete.Index, len(n.chain))
				return nil
			}

			blk := n.chain[concrete.Index]

			event = n.Send(&Packet{
				MessageType: msgtype.GetNewBlockRespMessageType,
				Payload:     &GetBlockDataResp{Block: blk},
				Source:      n,
				Destination: packet.Source,
				SizeInBytes: blockSize,
			}, now)
		}

		if n.serviceCode == servicecode.MinerNode {
			event = n.Send(&Packet{
				MessageType: msgtype.GetNewBlockRespMessageType,
				Payload:     &GetBlockDataResp{Block: n.newMinedBlock},
				Source:      n,
				Destination: packet.Source,
				SizeInBytes: blockSize,
			}, now)
		}

		return base.Aggregate(event)

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) getNewBlockDataRespHandler(packet *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *GetBlockDataResp:
		newBlock := concrete.Block

		//if newBlock.Index > n.inventory && newBlock == n.chain[len(n.chain)-1] {
		//	n.chain = MasterBlockchain
		//	n.inventory = newBlock.Index
		//
		//	return nil
		//}

		if newBlock.Index == n.chain[len(n.chain)-1].Index {
			return nil
		}

		//valid := IsBlockValid(newBlock, n.chain[len(n.chain)-1])
		//if !valid {
		//	fmt.Println("invalid block", newBlock.Index)
		//	event := n.Send(&Packet{
		//		MessageType: msgtype.NewBlockAckMessageType,
		//		Payload:     false,
		//		Source:      n,
		//		Destination: packet.Source,
		//	}, now)
		//
		//	//fmt.Println(n.name, "gets invalid block", newBlock.Index)
		//
		//	return base.Aggregate(event)
		//}

		n.chain = append(n.chain, newBlock)
		n.inventory = newBlock.Index

		if count == 13500 {
			fmt.Println("relay 90% nodes", now.String())
		}

		event := n.Send(&Packet{
			MessageType: msgtype.NewBlockAckMessageType,
			Payload:     true,
			Source:      n,
			Destination: packet.Source,
		}, now)

		events := n.sendPacketToPeerNodes(&Packet{
			MessageType: msgtype.InventoryMessage,
			Payload:     &InventoryMessage{Inventory: newBlock.Index},
			Source:      n,
		}, nodes, now)

		return append(events, event)

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) handleNewBlockAckMessage(packet *Packet, now time.Time) []base.Event {
	switch valid := packet.Payload.(type) {
	case bool:
		if valid {
		}

		return nil

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) handleGetBlockchain(request *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	event := n.Send(&Packet{
		MessageType: msgtype.GetBlockChainResp,
		Payload:     MasterBlockchain,
		Source:      n,
		Destination: request.Source}, now)

	return base.Aggregate(event)
}

func getLastMasterChainBlock() *Block {
	return MasterBlockchain[len(MasterBlockchain)-1]
}
