package bitcoin

import (
	"encoding/json"
	"p2psimulator/internal/bitcoin/dto"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/bitcoin/nodestate"
	"time"

	"github.com/bytedance/ns-x/v2/base"
)

func (n *Node) peerDiscoveryHandler(nodes map[string]base.Node) []base.Event {
	// Start peer discovery period
	var events []base.Event

	for _, dns := range n.DNSSeeds {
		destination := nodes[dns].(*Node)
		events = append(events, n.Send(&Packet{
			MessageType: msgtype.QueryMessageType,
			Payload:     []byte{},
			Source:      n,
			Destination: destination}, time.Now()),
		)
	}

	return events
}

func (n *Node) dnsAHandler(packet *Packet, nodes map[string]base.Node) []base.Event {
	var dnsAPayload dto.DNSARecord

	err := json.Unmarshal(packet.Payload, &dnsAPayload)
	if err != nil {
		n.logger.Error("failed unmarshal dnsA message")
		return nil
	}

	if dnsAPayload.IP != "" {
		n.State = nodestate.Connecting
	}

	n.availablePeers = append(n.availablePeers, dnsAPayload.IP)

	destination := nodes[dnsAPayload.IP].(*Node)
	n.Send(&Packet{
		MessageType: msgtype.VersionMessageType,
		Payload:     []byte{},
		Source:      n,
		Destination: destination}, time.Now())

	return nil
}
