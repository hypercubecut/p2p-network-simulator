package msgtype

type MessageType string

const (
	PingMessageType  MessageType = "ping"
	PongMessageType  MessageType = "pong"
	StartMessageType MessageType = "health_check"

	// bitcoin network message type
	GetBlockChainReq  MessageType = "get_block_chain_req"
	GetBlockChainResp MessageType = "get_block_chain_resp"

	WriteBlockReq  MessageType = "write_block_req"
	WriteBlockResp MessageType = "write_block_resp"

	PeerDiscoveryMessageType MessageType = "peer_discovery"
	QueryMessageType         MessageType = "query"
	DNSARecordMessageType    MessageType = "dns_a"

	VersionMessageType     MessageType = "version"
	VersionMessageBackType MessageType = "version_back"

	VerAckMessageType     MessageType = "version_ack"
	VerAckBackMessageType MessageType = "version_ack_back"

	GetAddressesMessageType     MessageType = "get_addr"
	GetAddressesRespMessageType MessageType = "get_addr_resp"

	GetBlocksMessageType     MessageType = "get_blocks"
	GetBlocksRespMessageType MessageType = "get_blocks_resp"

	GetBlockDataMessageType     MessageType = "get_block_data"
	GetBlockDataRespMessageType MessageType = "get_block_data_resp"

	ErrMessageType MessageType = "error"
)
