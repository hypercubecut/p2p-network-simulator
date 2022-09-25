package dto

type Peers struct {
	Peers []string `json:"peers"`
}

type Error struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type PeerDiscovery struct {
	MessageID string `json:"id"`
}

type DNSARecord struct {
	IP string `json:"ip"`
}

type Version struct {
	MessageID string `json:"id"`
}

type Verack struct {
}

type GetAddresses struct {
}
