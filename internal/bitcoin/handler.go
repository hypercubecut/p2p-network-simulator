package bitcoin

import (
	"errors"
	"fmt"
	"p2psimulator/internal/bitcoin/msgtype"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/bytedance/ns-x/v2/base"
	"go.uber.org/zap"
)

var (
	ErrUnknownPayload = errors.New("unknown request payload")
	ErrUnknownPacket  = errors.New("unknown packet type")
)

func (n *Node) Handler(nodes map[string]base.Node, logger *zap.Logger) func(packet base.Packet, now time.Time) []base.Event {
	return func(packet base.Packet, now time.Time) []base.Event {
		message, ok := packet.(*Packet)
		if !ok {
			n.logger.Error("failed to parse packet", zap.Error(ErrUnknownPacket))
			return nil
		}

		switch message.GetMessageType() {
		case msgtype.StartMessageType:
			return n.handleStart(message, nodes)

		case msgtype.GetBlockChainReq:
			return n.handleGetBlockchain(message, nodes)

		case msgtype.GetBlockChainResp:
			// Todo: implement here
			return nil

		case msgtype.WriteBlockReq:
			return n.handleWriteBlock(message, nodes)

		case msgtype.WriteBlockResp:
			// Todo: implement here
			return nil

		case msgtype.PingMessageType:
			return n.handlePing(message)

		case msgtype.PongMessageType:
			logger.Debug(fmt.Sprintf(n.name + " get pong from " + message.Source.name))
			return nil

		case msgtype.PeerDiscoveryMessageType:
			return n.peerDiscoveryHandler(nodes)

		case msgtype.QueryMessageType:
			return n.queryMessageHandler(message)

		case msgtype.DNSARecordMessageType:
			return n.dnsAHandler(message, nodes)

		case msgtype.VersionMessageType:
			return n.versionMsgHandler(message)

		case msgtype.VersionMessageBackType:
			return n.versionMsgBackHandler(message)

		case msgtype.VerAckMessageType:
			return n.verAckMessageHandler(message)

		case msgtype.VerAckBackMessageType:
			return n.verAckBackMessageHandler(message)

		case msgtype.GetAddressesMessageType:
			return n.getAddressHandler(message)

		case msgtype.GetAddressesRespMessageType:
			return n.getAddressesRespHandler(message, nodes)

		case msgtype.GetBlocksMessageType:
			return n.getBlocksHandler(message)

		case msgtype.GetBlocksRespMessageType:
			return n.getBlocksRespHandler(message)

		case msgtype.GetBlockDataMessageType:
			return n.getBlockDataHandler(message)

		case msgtype.GetBlockDataRespMessageType:
			return n.getBlockDataRespHandler(message)

		default:
			return nil
		}
	}
}

func (n *Node) handleStart(request *Packet, nodes map[string]base.Node) []base.Event {
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

func (n *Node) handleGetBlockchain(request *Packet, nodes map[string]base.Node) []base.Event {
	event := n.Send(&Packet{
		MessageType: msgtype.GetBlockChainResp,
		Payload:     MasterBlockchain,
		Source:      n,
		Destination: request.Source}, time.Now())

	return base.Aggregate(event)
}

func (n *Node) handleWriteBlock(request *Packet, nodes map[string]base.Node) []base.Event {
	reqDTO, ok := request.Payload.(*WriteBlockRequest)
	if !ok {
		n.logger.Error("handleWriteBlock failed unmarshal payload")
		return n.handleErrResp(msgtype.WriteBlockResp, ErrUnknownPayload, request)
	}

	newBlock, err := GenerateBlock(MasterBlockchain[len(MasterBlockchain)-1], reqDTO.BPM)
	if err != nil {
		n.logger.Error("failed to generateBlock", zap.Error(err))
		return n.handleErrResp(msgtype.WriteBlockResp, err, request)
	}

	if IsBlockValid(newBlock, MasterBlockchain[len(MasterBlockchain)-1]) {
		newBlockchain := append(MasterBlockchain, newBlock)
		ReplaceChain(newBlockchain)
		spew.Dump(MasterBlockchain)

		n.logger.Info("enter a new BPM")
	}

	respDTO := &WriteBlockResp{newBlock}

	event := n.Send(&Packet{
		MessageType: msgtype.WriteBlockResp,
		Payload:     respDTO,
		Source:      n,
		Destination: request.Source}, time.Now())

	return base.Aggregate(event)
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

func (n *Node) handlePing(request *Packet) []base.Event {
	from := request.Source

	n.logger.Info(fmt.Sprintf(n.ID()+" get ping message from %s", from.name))

	return base.Aggregate(n.Send(NewPacket(msgtype.PongMessageType, nil, n, from), time.Now()))
}

func (n *Node) handleOffline(request *Packet) []base.Event {
	from := request.Source

	n.logger.Error(fmt.Sprintf(n.ID() + " is offline cannot process request"))

	return base.Aggregate(n.Send(NewPacket(msgtype.PongMessageType,
		&Error{Msg: "server offline", Code: 500}, n, from), time.Now()))
}
