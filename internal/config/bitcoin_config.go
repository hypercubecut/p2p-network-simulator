package config

type bitcoinConfig struct {
	ServerToSeeds map[string][]string `json:"server_to_seeds"`
}
