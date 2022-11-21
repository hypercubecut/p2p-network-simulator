package bitcoin

import (
	"fmt"
	"p2psimulator/internal/bitcoin/msgtype"
	"time"

	"go.uber.org/zap"

	"github.com/bytedance/ns-x/v2/base"
)

func (n *Node) peerDiscoveryHandler(nodes map[string]base.Node, now time.Time) []base.Event {
	// Start peer discovery period
	var events []base.Event

	for _, dns := range n.seeds {
		destination := nodes[dns].(*Node)
		events = append(events, n.Send(&Packet{
			MessageType: msgtype.QueryMessageType,
			Payload:     []byte{},
			Source:      n,
			Destination: destination}, now),
		)
		n.logger.Debug(fmt.Sprintf("%s send query message to %s", n.name, destination.name))
	}

	return events
}

func (n *Node) queryMessageHandler(packet *Packet, now time.Time) []base.Event {
	respPayload := &DNSARecord{IP: n.name}

	event := n.Send(&Packet{
		MessageType: msgtype.DNSARecordMessageType,
		Payload:     respPayload,
		Source:      n,
		Destination: packet.Source}, now)

	n.logger.Debug(fmt.Sprintf("%s send DNS A record message to %s", n.name, packet.Source.name))

	return base.Aggregate(event)
}

func (n *Node) dnsAHandler(packet *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	dnsAPayload, ok := packet.Payload.(*DNSARecord)
	if !ok {
		n.logger.Error("dnsAHandler failed unmarshal dnsA message")
		return n.handleErrResp(msgtype.VersionMessageType, ErrUnknownPayload, packet)
	}

	n.logger.Debug(fmt.Sprintf("%s get DNS A record message from %s, %s",
		n.name, packet.Source.name, dnsAPayload.IP))

	event := n.Send(&Packet{
		MessageType: msgtype.VersionMessageType,
		Payload: &VersionMessage{
			Version:   defaultVersion,
			Services:  n.serviceCode,
			Timestamp: now.UnixMilli(),
		},
		Source:      n,
		Destination: packet.Source}, now)

	return base.Aggregate(event)
}

func (n *Node) versionMsgHandler(packet *Packet, now time.Time) []base.Event {
	switch packet.Payload.(type) {
	case *VersionMessage:
		_ = packet.Payload.(*VersionMessage)
	case *Error:
		n.logger.Error("versionMsg is error",
			zap.String("err", packet.Payload.(*Error).Msg))

		return n.handleErrResp(msgtype.VersionMessageBackType, ErrUnknownPayload, packet)
	}

	// Todo: validate version msg here

	event := n.Send(&Packet{
		MessageType: msgtype.VersionMessageBackType,
		Payload: &VersionMessage{
			Version:   defaultVersion,
			Services:  n.serviceCode,
			Timestamp: time.Now().UnixMilli(),
		},
		Source:      n,
		Destination: packet.Source}, now)

	return base.Aggregate(event)
}

func (n *Node) versionMsgBackHandler(packet *Packet, now time.Time) []base.Event {
	switch packet.Payload.(type) {
	case *VersionMessage:
		_ = packet.Payload.(*VersionMessage)
	case *Error:
		n.logger.Error("versionMsgBack is error",
			zap.String("err", packet.Payload.(*Error).Msg))

		return n.handleErrResp(msgtype.VerAckMessageType, ErrUnknownPayload, packet)
	}

	event := n.Send(&Packet{
		MessageType: msgtype.VerAckMessageType,
		Payload:     &VersionAckMessage{},
		Source:      n,
		Destination: packet.Source}, now)

	return base.Aggregate(event)
}

func (n *Node) verAckMessageHandler(packet *Packet, now time.Time) []base.Event {
	switch packet.Payload.(type) {
	case *VersionAckMessage:
		_ = packet.Payload.(*VersionAckMessage)
	case *Error:
		n.logger.Error("verAckMessage is error",
			zap.String("err", packet.Payload.(*Error).Msg))

		return n.handleErrResp(msgtype.VerAckBackMessageType, ErrUnknownPayload, packet)
	}

	n.AddNewPeers(packet.Source.name)

	n.logger.Debug(fmt.Sprintf("%s get ver ack message from %s", n.name, packet.Source.name))

	event := n.Send(&Packet{
		MessageType: msgtype.VerAckBackMessageType,
		Payload:     &VersionAckMessage{},
		Source:      n,
		Destination: packet.Source}, now)

	return base.Aggregate(event)
}

func (n *Node) verAckBackMessageHandler(packet *Packet, now time.Time) []base.Event {
	switch packet.Payload.(type) {
	case *VersionAckMessage:
		n.AddNewPeers(packet.Source.name)
		n.logger.Debug(fmt.Sprintf("%s get ver ack back message from %s", n.name, packet.Source.name))

		event := n.Send(&Packet{
			MessageType: msgtype.GetAddressesMessageType,
			Payload:     nil,
			Source:      n,
			Destination: packet.Source,
		}, now)

		return base.Aggregate(event)

	case *Error:
		n.logger.Error("verAckBackMessage is error",
			zap.String("err", packet.Payload.(*Error).Msg))
		return nil

	default:
		return nil
	}
}

func (n *Node) getAddressHandler(packet *Packet, now time.Time) []base.Event {
	respDTO := &GetAddressResp{
		MorePeers: n.GetAvailablePeers(),
	}

	event := n.Send(&Packet{
		MessageType: msgtype.GetAddressesRespMessageType,
		Payload:     respDTO,
		Source:      n,
		Destination: packet.Source,
	}, now)

	return base.Aggregate(event)
}

func (n *Node) getAddressesRespHandler(packet *Packet, nodes map[string]base.Node, now time.Time) []base.Event {
	switch packet.Payload.(type) {
	case *GetAddressResp:
		n.AddNewPeers(packet.Payload.(*GetAddressResp).MorePeers...)

		// Todo: implement header first here
		if n.state == "IBD" {
			return nil
		}

		n.logger.Debug(fmt.Sprintf("start initial block download for node %s", n.name))
		n.state = "IBD"

		return n.initialBlockDownloadWithBlocksFirst(nodes, now)

	default:
		return nil
	}
}
