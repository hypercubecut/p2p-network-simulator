package bitcoin

import (
	"errors"
	"p2psimulator/internal/bitcoin/msgtype"
	"time"

	"github.com/bytedance/ns-x/v2/base"
)

const (
	maxInventoryLen = 500
)

func (n *Node) initialBlockDownload(nodes map[string]base.Node) []base.Event {
	// pick one full node
	var pickedNode *Node
	for peer, _ := range n.availablePeers {
		pn, ok := nodes[peer].(*Node)
		if !ok {
			return nil
		}

		if pn.GetServiceCode() == NODE_NETWORK {
			pickedNode = pn
			break
		}
	}

	if pickedNode == nil {
		return nil
	}

	event := n.Send(&Packet{
		MessageType: msgtype.GetBlocksMessageType,
		Payload: &GetBlocksReq{
			Version:    0,
			BlockIndex: n.inventory,
		},
		Source:      n,
		Destination: pickedNode,
	}, time.Now())

	return base.Aggregate(event)
}

func (n *Node) getBlocksHandler(packet *Packet) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *GetBlocksReq:
		clientLastBlockIdx := concrete.BlockIndex

		var inv int
		if n.inventory > clientLastBlockIdx {
			if n.inventory-clientLastBlockIdx > maxInventoryLen {
				inv = clientLastBlockIdx + 1 + maxInventoryLen
			} else {
				inv = n.inventory
			}
		}

		event := n.Send(&Packet{
			Packet:      nil,
			MessageType: msgtype.GetBlocksRespMessageType,
			Payload: &GetBlocksResp{
				Version:   n.version,
				Inventory: inv,
			},
			Source:      n,
			Destination: packet.Source,
		}, time.Now())

		return base.Aggregate(event)

	default:
		return n.handleErrResp(msgtype.GetBlocksRespMessageType, ErrUnknownPayload, packet)
	}
}

func (n *Node) getBlocksRespHandler(packet *Packet) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *GetBlocksResp:
		inv := concrete.Inventory
		if inv == 0 {
			return nil
		}

		if inv-n.inventory >= maxInventoryLen {
			n.inventory = inv

			event := n.Send(&Packet{
				MessageType: msgtype.GetBlocksMessageType,
				Payload: &GetBlocksReq{
					Version:    n.version,
					BlockIndex: inv,
				},
				Source:      n,
				Destination: packet.Source,
			}, time.Now())

			return base.Aggregate(event)
		}

		n.inventory = inv

		return n.getData(packet)

	case *Error:
		return nil

	default:
		return nil
	}
}

func (n *Node) getData(packet *Packet) []base.Event {
	// get missing blocks' data
	if n.inventory > n.chain[len(n.chain)-1].Index {
		event := n.Send(&Packet{
			MessageType: msgtype.GetBlockDataMessageType,
			Payload:     &GetBlockDataReq{n.chain[len(n.chain)-1].Index + 1},
			Source:      n,
			Destination: packet.Source,
		}, time.Now())

		return base.Aggregate(event)
	}

	return nil
}

func (n *Node) getBlockDataHandler(packet *Packet) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *GetBlockDataReq:
		if concrete.Index > len(n.chain) || n.chain[concrete.Index] == nil {
			return n.handleErrResp(msgtype.GetBlockDataRespMessageType,
				errors.New("node does not have the block"), packet)
		}

		event := n.Send(&Packet{
			MessageType: msgtype.GetBlockDataRespMessageType,
			Payload:     &GetBlockDataResp{n.chain[concrete.Index]},
			Source:      n,
			Destination: packet.Source,
		}, time.Now())

		return base.Aggregate(event)

	default:
		return n.handleErrResp(msgtype.GetBlockDataRespMessageType, ErrUnknownPayload, packet)
	}
}

func (n *Node) getBlockDataRespHandler(packet *Packet) []base.Event {
	switch concrete := packet.Payload.(type) {
	case *GetBlockDataResp:
		blk := concrete.Block
		n.chain = append(n.chain, blk)

		if n.inventory > n.chain[len(n.chain)-1].Index {
			event := n.Send(&Packet{
				MessageType: msgtype.GetBlockDataMessageType,
				Payload:     &GetBlockDataReq{n.chain[len(n.chain)-1].Index + 1},
				Source:      n,
				Destination: packet.Source,
			}, time.Now())

			return base.Aggregate(event)
		}

		return nil

	default:
		// Todo: do something here to get the data
		return nil
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
