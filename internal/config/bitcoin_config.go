package config

type BitcoinConfig struct {
	MasterChainLen int `json:"master_chain_len"`
	NumberOfPeers  int `json:"number_of_peers"`
}
