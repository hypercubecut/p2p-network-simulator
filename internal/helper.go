package internal

import (
	"p2psimulator/internal/bitcoin/servicecode"
	"p2psimulator/internal/config"
	"strconv"
)

type Option func(cfg *config.Config)

func apply(cfg *config.Config, options ...Option) {
	for _, option := range options {
		option(cfg)
	}
}

func GenerateConfig(initialChainLen, numOfFullNodes, numOfNewNodes, numberOfMiner int, options ...Option) *config.Config {
	cfg := &config.Config{
		SimulatorCfg: &config.SimulatorConfig{
			EnableDebugLog: false,
			LifeTimeInMin:  12000,
		},
		ServersCfg: &config.ServersConfig{
			Servers: []config.ServerDetailsConfig{},
		},
		BitcoinCfg: &config.BitcoinConfig{MasterChainLen: initialChainLen},
	}

	idx := 1
	var allFullNodes []string
	for idx < numOfFullNodes+1 {
		name := "f" + strconv.Itoa(idx)
		allFullNodes = append(allFullNodes, name)
		cfg.ServersCfg.Servers = append(cfg.ServersCfg.Servers,
			config.ServerDetailsConfig{
				Name:            name,
				ServiceCode:     servicecode.FullNode,
				Version:         7000,
				Seeds:           nil,
				OutputDelayInMs: 200,
				OutputBPS:       -1,
				OutputPPS:       10000,
				InputDelayInMs:  200,
				InputBPS:        -1,
				InputPPS:        10000,
				QueueLimit:      10000,
			})

		idx++
	}

	cfg.ServersCfg.AllFullNodes = allFullNodes

	var allNewNodes []string
	idx = 1
	for idx < numOfNewNodes+1 {
		allNewNodes = append(allNewNodes, "p"+strconv.Itoa(idx))
		cfg.ServersCfg.Servers = append(cfg.ServersCfg.Servers,
			config.ServerDetailsConfig{
				Name:            "p" + strconv.Itoa(idx),
				ServiceCode:     servicecode.Unnamed,
				Seeds:           allFullNodes,
				Version:         7000,
				OutputDelayInMs: 200,
				OutputBPS:       -1,
				OutputPPS:       1000,
				InputDelayInMs:  200,
				InputBPS:        -1,
				InputPPS:        1000,
				QueueLimit:      1000,
			})

		idx++
	}

	cfg.ServersCfg.AllNewNodes = allNewNodes

	var allMinerNodes []string
	idx = 1
	for idx < numberOfMiner+1 {
		allNewNodes = append(allNewNodes, "m"+strconv.Itoa(idx))
		cfg.ServersCfg.Servers = append(cfg.ServersCfg.Servers,
			config.ServerDetailsConfig{
				Name:            "m" + strconv.Itoa(idx),
				ServiceCode:     servicecode.MinerNode,
				Seeds:           allFullNodes,
				Version:         7000,
				OutputDelayInMs: 200,
				OutputBPS:       -1,
				OutputPPS:       1000,
				InputDelayInMs:  200,
				InputBPS:        -1,
				InputPPS:        1000,
				QueueLimit:      1000,
			})

		idx++
	}

	cfg.ServersCfg.AllMinersNodes = allMinerNodes

	apply(cfg, options...)

	return cfg
}

func WithFullNodeDelay(delayInMs int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.FullNode {
				server.InputDelayInMs = delayInMs
				server.OutputDelayInMs = delayInMs
			}
		}
	}
}

func WithFullNodePPS(pps int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.FullNode {
				server.InputPPS = pps
				server.OutputPPS = pps
			}
		}
	}
}

func WithFullNodeBPS(bps int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.FullNode {
				server.InputBPS = bps
				server.InputBPS = bps
			}
		}
	}
}

func WithNewNodeDelay(delayInMs int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.Unnamed {
				server.InputDelayInMs = delayInMs
				server.OutputDelayInMs = delayInMs
			}
		}
	}
}

func WithNewNodePPS(pps int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.Unnamed {
				server.InputPPS = pps
				server.OutputPPS = pps
			}
		}
	}
}

func WithNewNodeBPS(bps int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.Unnamed {
				server.InputBPS = bps
				server.InputBPS = bps
			}
		}
	}
}

func WithQueueLimit(limit int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			server.QueueLimit = limit
		}
	}
}

func WithMinerNodeDelay(delayInMs int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.MinerNode {
				server.InputDelayInMs = delayInMs
				server.OutputDelayInMs = delayInMs
			}
		}
	}
}

func WithMinerNodePPS(pps int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.MinerNode {
				server.InputPPS = pps
				server.OutputPPS = pps
			}
		}
	}
}

func WithMinerNodeBPS(bps int64) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.MinerNode {
				server.InputBPS = bps
				server.InputBPS = bps
			}
		}
	}
}
