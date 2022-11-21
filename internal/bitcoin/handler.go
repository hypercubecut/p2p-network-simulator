package bitcoin

import (
	"errors"
	"fmt"
	"p2psimulator/internal/bitcoin/msgtype"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"go.uber.org/zap"
)

var (
	ErrUnknownPayload = errors.New("unknown request payload")
	ErrUnknownPacket  = errors.New("unknown packet type")
)

func (n *Node) Handler(nodes map[string]base.Node, now time.Time, logger *zap.Logger) func(packet base.Packet, now time.Time) []base.Event {
	return func(packet base.Packet, now time.Time) []base.Event {
		message, ok := packet.(*Packet)
		if !ok {
			n.logger.Error("failed to parse packet", zap.Error(ErrUnknownPacket))
			return nil
		}

		switch message.GetMessageType() {
		case msgtype.StartMessageType:
			return n.handleStart(message, nodes, now)

		case msgtype.GetBlockChainReq:
			return n.handleGetBlockchain(message, nodes, now)

		case msgtype.GetBlockChainResp:
			return nil

		case msgtype.MineNewBlockReq:
			return n.handleMineNewBlock(message, nodes, now)

		case msgtype.InventoryMessage:
			return n.handleInventoryMessage(message, now)

		case msgtype.InvalidInventoryMessage:
			return n.handleInvalidInventory(nodes, now)

		case msgtype.MineNewBlockResp:
			return nil

		case msgtype.PingMessageType:
			return n.handlePing(message, now)

		case msgtype.PongMessageType:
			logger.Debug(fmt.Sprintf(n.name + " get pong from " + message.Source.name))
			return nil

		case msgtype.PeerDiscoveryMessageType:
			return n.peerDiscoveryHandler(nodes, now)

		case msgtype.QueryMessageType:
			return n.queryMessageHandler(message, now)

		case msgtype.DNSARecordMessageType:
			return n.dnsAHandler(message, nodes, now)

		case msgtype.VersionMessageType:
			return n.versionMsgHandler(message, now)

		case msgtype.VersionMessageBackType:
			return n.versionMsgBackHandler(message, now)

		case msgtype.VerAckMessageType:
			return n.verAckMessageHandler(message, now)

		case msgtype.VerAckBackMessageType:
			return n.verAckBackMessageHandler(message, now)

		case msgtype.GetAddressesMessageType:
			return n.getAddressHandler(message, now)

		case msgtype.GetAddressesRespMessageType:
			return n.getAddressesRespHandler(message, nodes, now)

		case msgtype.GetBlocksMessageType:
			return n.getBlocksHandler(message, now)

		case msgtype.GetBlocksRespMessageType:
			return n.getBlocksRespHandler(message, now)

		case msgtype.GetBlockDataMessageType:
			return n.getBlockDataHandler(message, now)

		case msgtype.GetBlockDataRespMessageType:
			return n.getBlockDataRespHandler(message, now)

		case msgtype.GetNewBlockMessageType:
			return n.getNewBlockDataHandler(message, now)

		case msgtype.GetNewBlockRespMessageType:
			return n.getNewBlockDataRespHandler(message, nodes, now)

		case msgtype.NewBlockAckMessageType:
			return n.handleNewBlockAckMessage(message, now)

		case msgtype.P2PWithBroadcastMessageType:
			return n.handleP2PWithBroadcast(message, nodes, now)

		default:
			return nil
		}
	}
}

func (n *Node) handleStart(request *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	var reqDTO *Peers

	switch request.Payload.(type) {
	case *Peers:
		reqDTO = request.Payload.(*Peers)

	case *Error:
		return nil

	default:
		n.logger.Error(fmt.Sprintf("handleStart failed unmarshal payload %+v", request.Payload),
			zap.Error(ErrUnknownPayload))
		return nil
	}

	n.seeds = reqDTO.Peers

	n.logger.Info(fmt.Sprintf(n.name+" receive start message at %s", time.Now().String()))
	n.logger.Info(fmt.Sprintf(n.name+" start sending ping to all seeds at %s", time.Now().String()))

	var events []base.Event
	for _, dns := range n.seeds {
		destination := nodes[dns].(*Node)

		events = append(events, n.Send(&Packet{
			MessageType: msgtype.PingMessageType,
			Payload:     []byte{},
			Source:      n,
			Destination: destination}, time.Now()),
		)
	}

	// send ping to peer
	return events
}

func (n *Node) handleErrResp(respMsgType msgtype.MessageType, err error, request *Packet) []base.Event {
	if request.Source == nil {
		return nil
	}

	event := n.Send(&Packet{
		MessageType: respMsgType,
		Payload: &Error{
			Msg:  err.Error(),
			Code: 0,
		},
		Source:      n,
		Destination: request.Source}, time.Now())

	return base.Aggregate(event)
}

func (n *Node) handlePing(request *Packet, now time.Time) []base.Event {
	from := request.Source

	n.logger.Info(fmt.Sprintf(n.ID()+" get ping message from %s", from.name))

	return base.Aggregate(n.Send(NewPacket(msgtype.PongMessageType, nil, n, from), now))
}

func (n *Node) handleOffline(request *Packet, now time.Time) []base.Event {
	from := request.Source

	n.logger.Error(fmt.Sprintf(n.ID() + " is offline cannot process request"))

	return base.Aggregate(n.Send(NewPacket(msgtype.PongMessageType,
		&Error{Msg: "server offline", Code: 500}, n, from), now))
}
