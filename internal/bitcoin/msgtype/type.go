package msgtype

type MessageType string

const (
	PingMessageType  MessageType = "ping"
	PongMessageType  MessageType = "pong"
	StartMessageType MessageType = "health_check"

	// bitcoin network message type
	PeerDiscoveryMessageType MessageType = "peer_discovery"
	QueryMessageType         MessageType = "query"
	DNSARecordMessageType    MessageType = "dns_a"
	VersionMessageType       MessageType = "version"
	VerackMessageType        MessageType = "verack"
	GetAddressesMessageType  MessageType = "addr"
	ErrMessageType           MessageType = "error"
)
