package bitcoin

import (
	"fmt"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/bitcoin/servicecode"
	"time"

	"go.uber.org/zap"

	"github.com/bytedance/ns-x/v2/base"
)

func (n *Node) handleMineNewBlock(request *Packet, nodes map[string]base.Node) []base.Event {
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

	n.newMinedBlock = newBlock
	n.logger.Info(fmt.Sprintf("%s finish mine a new block", n.name), zap.Int("newBlock", newBlock.Index))

	return n.handleStandardBlockRelay(newBlock, nodes)
}

func (n *Node) handleStandardBlockRelay(newBlock *Block, nodes map[string]base.Node) []base.Event {
	var events []base.Event

	for _, node := range nodes {
		dest, ok := node.(*Node)
		if ok {
			if dest.name == n.name || dest.serviceCode != servicecode.FullNode {
				continue
			}

			event := n.Send(&Packet{
				MessageType: msgtype.InventoryMessage,
				Payload:     &InventoryMessage{Inventory: newBlock.Index},
				Source:      n,
				Destination: dest,
			}, time.Now())

			events = append(events, event)
		}
	}

	return events
}

func (n *Node) handleInventoryMessage(packet *Packet, nodes map[string]base.Node) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *InventoryMessage:
		// full node request for the new block if valid
		if n.serviceCode == servicecode.FullNode && n.inventory+1 == concrete.Inventory {
			var events []base.Event

			events = append(events, n.Send(&Packet{
				MessageType: msgtype.GetNewBlockMessageType,
				Payload: &GetBlockDataReq{
					Index: concrete.Inventory,
				},
				Source:      n,
				Destination: packet.Source,
			}, time.Now()))

			return events
		}

		n.logger.Info(fmt.Sprintf("%s received invalid/notrequired inventory message from %s",
			n.name, packet.Source.name),
			zap.Int("new_inventory", concrete.Inventory),
			zap.Int("my_inventory", n.GetInventory()))

		return nil

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) getNewBlockDataHandler(packet *Packet) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *GetBlockDataReq:
		if n.serviceCode != servicecode.MinerNode ||
			n.newMinedBlock == nil || concrete.Index != n.newMinedBlock.Index {
			return nil
		}

		event := n.Send(&Packet{
			MessageType: msgtype.GetNewBlockRespMessageType,
			Payload:     &GetBlockDataResp{Block: n.newMinedBlock},
			Source:      n,
			Destination: packet.Source,
		}, time.Now())

		return base.Aggregate(event)

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) getNewBlockDataRespHandler(packet *Packet) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *GetBlockDataResp:
		newBlock := concrete.Block

		mclock.Lock()
		defer mclock.Unlock()

		if newBlock.Index > n.inventory && newBlock == getLastMasterChainBlock() {
			n.chain = MasterBlockchain
			n.inventory = newBlock.Index

			return nil
		}

		valid := IsBlockValid(newBlock, getLastMasterChainBlock())
		if !valid {
			event := n.Send(&Packet{
				MessageType: msgtype.NewBlockAckMessageType,
				Payload:     false,
				Source:      n,
				Destination: packet.Source,
			}, time.Now())

			return base.Aggregate(event)
		}

		MasterBlockchain = append(n.chain, newBlock)
		n.chain = MasterBlockchain
		n.inventory = newBlock.Index

		event := n.Send(&Packet{
			MessageType: msgtype.NewBlockAckMessageType,
			Payload:     true,
			Source:      n,
			Destination: packet.Source,
		}, time.Now())

		return base.Aggregate(event)

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) handleNewBlockAckMessage(packet *Packet) []base.Event {
	switch valid := packet.Payload.(type) {
	case bool:
		if valid && n.inventory < n.newMinedBlock.Index {
			n.inventory = n.newMinedBlock.Index
		}

		return nil

	default:
		n.logger.Error("failed to convert payload", zap.Error(ErrUnknownPayload))

		return nil
	}
}

func (n *Node) handleGetBlockchain(request *Packet, nodes map[string]base.Node) []base.Event {
	event := n.Send(&Packet{
		MessageType: msgtype.GetBlockChainResp,
		Payload:     MasterBlockchain,
		Source:      n,
		Destination: request.Source}, time.Now())

	return base.Aggregate(event)
}

func getLastMasterChainBlock() *Block {
	return MasterBlockchain[len(MasterBlockchain)-1]
}
