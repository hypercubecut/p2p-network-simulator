package bitcoin

type WriteBlockRequest struct {
	BPM int
}

type WriteBlockResp struct {
	*Block
}

type VersionMessage struct {
	Version   int32
	Services  uint64
	Timestamp int64
}

type VersionAckMessage struct {
}

type Peers struct {
	Peers []string `json:"peers"`
}

type Error struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type DNSARecord struct {
	IP string `json:"ip"`
}

type GetAddressResp struct {
	MorePeers []string
}

type GetBlocksReq struct {
	Version    int32
	BlockIndex int
}

type GetBlocksResp struct {
	Version   int32
	Inventory int
}

type GetBlockDataReq struct {
	Index int
}

type GetBlockDataResp struct {
	Block *Block
}
