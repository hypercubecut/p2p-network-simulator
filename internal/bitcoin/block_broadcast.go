package bitcoin

import (
	"fmt"
	"p2psimulator/internal/bitcoin/msgtype"
	"time"

	"go.uber.org/zap"

	"github.com/bytedance/ns-x/v2/base"
)

func (n *Node) handleMineNewBlock(request *Packet, nodes map[string]base.Node) []base.Event {
	n.logger.Debug(fmt.Sprintf("%s start mine a new block", n.name), zap.Int("chainLen", len(n.chain)))

	reqDTO, ok := request.Payload.(*WriteBlockRequest)
	if !ok {
		n.logger.Error("handleWriteBlock failed unmarshal payload")
		return n.handleErrResp(msgtype.MineNewBlockResp, ErrUnknownPayload, request)
	}

	newBlock, err := GenerateBlock(n.chain[len(n.chain)-1], reqDTO.BPM)
	if err != nil {
		n.logger.Error(fmt.Sprintf("%s failed to generateBlock", n.name), zap.Error(err))
		return n.handleErrResp(msgtype.MineNewBlockResp, err, request)
	}

	if IsBlockValid(newBlock, n.chain[len(n.chain)-1]) {
		newBlockchain := append(n.chain, newBlock)
		n.ReplaceChain(newBlockchain)
		//spew.Dump(n.chain)

		n.logger.Debug(fmt.Sprintf("%s finish mine a new block", n.name), zap.Int("chainLen", len(n.chain)))

		n.inventory = newBlock.Index

		return n.handleStandardBlockRelay(newBlock, nodes)
	}

	return nil
}

func (n *Node) handleStandardBlockRelay(newBlock *Block, nodes map[string]base.Node) []base.Event {
	var events []base.Event

	for peer, ok := range n.availablePeers {
		if ok {
			dest, ok := nodes[peer].(*Node)
			if !ok {
				return nil
			}

			if dest.name == n.name {
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

func (n *Node) handleInventoryMessage(packet *Packet) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *InventoryMessage:
		// full node request for the new block if valide
		if n.inventory+1 == concrete.Inventory {
			event := n.Send(&Packet{
				MessageType: msgtype.GetBlockDataMessageType,
				Payload: &GetBlockDataReq{
					Index: concrete.Inventory,
				},
				Source:      n,
				Destination: packet.Source,
			}, time.Now())

			n.inventory = concrete.Inventory

			n.logger.Debug(fmt.Sprintf("%s received valid inventory message from %s",
				n.name, packet.Source.name), zap.Int("inventory", concrete.Inventory))

			return base.Aggregate(event)
		}

		n.logger.Debug(fmt.Sprintf("%s received invalid/notrequired inventory message from %s",
			n.name, packet.Source.name), zap.Int("inventory", concrete.Inventory))

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
