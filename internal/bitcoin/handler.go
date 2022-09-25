package bitcoin

import (
	"encoding/json"
	"fmt"
	"p2psimulator/internal/bitcoin/dto"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/bitcoin/nodestate"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"go.uber.org/zap"
)

func (n *Node) Handler(nodes map[string]base.Node, logger *zap.Logger) func(packet base.Packet, now time.Time) []base.Event {
	return func(packet base.Packet, now time.Time) []base.Event {
		message, ok := packet.(*Packet)
		if !ok {
			n.logger.Error("failed to parse packet")

			return nil
		}

		if n.State == nodestate.Offline {
			return n.handleOffline(message)
		}

		switch message.GetMessageType() {
		case msgtype.StartMessageType:
			return n.handleStart(message, nodes)

		case msgtype.PingMessageType:
			return n.handlePing(message)

		case msgtype.PongMessageType:
			logger.Info(fmt.Sprintf(n.name + " get pong from " + message.Source.name))
			return nil

		case msgtype.PeerDiscoveryMessageType:
			return n.peerDiscoveryHandler(nodes)

		case msgtype.DNSARecordMessageType:
			return n.dnsAHandler(message, nodes)

		case msgtype.QueryMessageType:
			payload, _ := json.Marshal(&dto.DNSARecord{IP: "192.0.2.113"})

			event := n.Send(&Packet{
				MessageType: msgtype.DNSARecordMessageType,
				Payload:     payload,
				Source:      n,
				Destination: message.Source}, now)

			return base.Aggregate(event)

		default:
			return nil
		}
	}
}

func (n *Node) handleStart(request *Packet, nodes map[string]base.Node) []base.Event {
	var seeds dto.Peers

	err := json.Unmarshal(request.Payload, &seeds)
	if err != nil {
		n.logger.Error("failed unmarshal payload",
			zap.String("payload", string(request.Payload)),
			zap.Error(err))
		return nil
	}

	n.DNSSeeds = seeds.Peers

	n.logger.Info(fmt.Sprintf(n.name+" receive start message at %s", time.Now().String()))
	n.logger.Info(fmt.Sprintf(n.name+" start sending ping to all seeds at %s", time.Now().String()))

	var events []base.Event
	for _, dns := range n.DNSSeeds {
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

func (n *Node) handlePing(request *Packet) []base.Event {
	from := request.Source

	n.logger.Info(fmt.Sprintf(n.ID()+" get ping message from %s", from.name))

	return base.Aggregate(n.Send(NewPacket(msgtype.PongMessageType, nil, n, from), time.Now()))
}

func (n *Node) handleOffline(request *Packet) []base.Event {
	from := request.Source

	n.logger.Error(fmt.Sprintf(n.ID() + " is offline cannot process request"))

	return base.Aggregate(n.Send(NewPacket(msgtype.PongMessageType,
		&dto.Error{Msg: "server offline", Code: 500}, n, from), time.Now()))
}
