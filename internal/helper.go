package internal

import (
	"math/rand"
	"p2psimulator/internal/bitcoin/servicecode"
	"p2psimulator/internal/config"
	"strconv"
	"time"
)

type Option func(cfg *config.Config)

func apply(cfg *config.Config, options ...Option) {
	for _, option := range options {
		option(cfg)
	}
}

func GenerateConfig(
	initialChainLen, numOfFullNodes, numOfNewNodes, numberOfMiner, numberOfPeers, peerOp int, options ...Option) *config.Config {
	cfg := &config.Config{
		SimulatorCfg: &config.SimulatorConfig{
			EnableDebugLog: false,
			LifeTimeInMin:  1000,
		},
		ServersCfg: &config.ServersConfig{
			Servers: []*config.ServerDetailsConfig{},
		},
		BitcoinCfg: &config.BitcoinConfig{
			MasterChainLen: initialChainLen,
			NumberOfPeers:  numberOfPeers,
		},
	}

	getPeersFunc := getPeers(peerOp)

	idx := 1
	count := 0
	var allFullNodes []string
	for count < numOfFullNodes {
		name := "p-" + strconv.Itoa(idx)
		allFullNodes = append(allFullNodes, name)
		cfg.ServersCfg.Servers = append(cfg.ServersCfg.Servers,
			&config.ServerDetailsConfig{
				Name:            name,
				ServiceCode:     servicecode.FullNode,
				Version:         7000,
				Seeds:           nil,
				OutputDelayInMs: 200,
				OutputBPS:       -1,
				InputDelayInMs:  200,
				InputBPS:        -1,
				QueueLimit:      10000,
				//Peers:           getNRandomPeers(numberOfPeers, 1, numOfFullNodes),
				Peers: getPeersFunc(numberOfPeers, idx, 1, numOfFullNodes),
			})

		idx++
		count++
	}

	cfg.ServersCfg.AllFullNodes = allFullNodes

	var allNewNodes []string
	count = 0
	for count < numOfNewNodes {
		allNewNodes = append(allNewNodes, "p-"+strconv.Itoa(idx))
		cfg.ServersCfg.Servers = append(cfg.ServersCfg.Servers,
			&config.ServerDetailsConfig{
				Name:            "p-" + strconv.Itoa(idx),
				ServiceCode:     servicecode.Unnamed,
				Seeds:           allFullNodes,
				Version:         7000,
				OutputDelayInMs: 200,
				OutputBPS:       -1,
				InputDelayInMs:  200,
				InputBPS:        -1,
				QueueLimit:      1000,
				//Peers:           getNRandomPeers(numberOfPeers, 1, numOfFullNodes),
				Peers: getPeersFunc(numberOfPeers, idx, 1, numOfFullNodes),
			})

		idx++
		count++
	}

	cfg.ServersCfg.AllNewNodes = allNewNodes

	var allMinerNodes []string
	count = 0
	for count < numberOfMiner {
		allMinerNodes = append(allMinerNodes, "p-"+strconv.Itoa(idx))
		cfg.ServersCfg.Servers = append(cfg.ServersCfg.Servers,
			&config.ServerDetailsConfig{
				Name:            "p-" + strconv.Itoa(idx),
				ServiceCode:     servicecode.MinerNode,
				Version:         7000,
				OutputDelayInMs: 200,
				OutputBPS:       -1,
				InputDelayInMs:  200,
				InputBPS:        -1,
				QueueLimit:      1000,
				//Peers:           getNRandomPeers(numberOfPeers, 1, numOfFullNodes),
				Peers: getPeersFunc(numberOfPeers, idx, 1, numOfFullNodes),
			})

		idx++
		count++
	}

	cfg.ServersCfg.AllMinersNodes = allMinerNodes

	apply(cfg, options...)

	//fmt.Println(fmt.Sprintf("%+v", cfg.ServersCfg.Servers[1500]))

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
				server.OutputBPS = bps
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
				server.OutputBPS = bps
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

func WithRandomMinerNodeDelay() Option {

	return func(cfg *config.Config) {
		rand.Seed(time.Now().UnixNano())
		min := 0
		max := 1000

		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.MinerNode {
				delta := rand.Intn(max-min+1) + min
				server.InputDelayInMs = int64(50 + delta)
				server.OutputDelayInMs = int64(50 + delta)
			}
		}
	}
}

func WithDifferentMinerNodeDelay(minerName string) Option {
	return func(cfg *config.Config) {
		for _, server := range cfg.ServersCfg.Servers {
			if server.ServiceCode == servicecode.MinerNode {
				if server.Name == minerName {
					server.InputDelayInMs = 10
					server.OutputDelayInMs = 10
				} else {
					server.InputDelayInMs = 200
					server.OutputDelayInMs = 200
				}
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
				server.OutputBPS = bps
			}
		}
	}
}

func getPeers(op int) func(num, my, low, high int) []string {
	if op == 1 {
		return getNClosestPeers
	}

	if op == 2 {
		return getNSeparatedPeers
	}

	return getNRandomPeers
}

func getNClosestPeers(num, my, low, high int) []string {
	rand.Seed(time.Now().UnixNano())

	var result []string

	if num == 1 {
		if my+1 <= high {
			return []string{"p-" + strconv.Itoa(my+1)}
		} else {
			return []string{"p-1"}
		}

	}

	lower := my - num/2
	higher := my + num/2 + 1

	if lower < low {
		higher += low - lower
		lower = low
	}

	if higher > high {
		lower -= higher - high
		higher = high
	}

	for i := lower; i <= higher; i++ {
		if i != my {
			result = append(result, "p-"+strconv.Itoa(i))
		}
	}

	return result
}

func getNRandomPeers(num, my, low, high int) []string {
	rand.Seed(time.Now().UnixNano())

	var result []string

	for num > 0 {
		randomNum := low + rand.Intn(high-low+1)
		result = append(result, "p-"+strconv.Itoa(randomNum))
		num--
	}

	return result
}

func getNSeparatedPeers(num, my, low, high int) []string {
	var result []string

	step := (high - low) / num

	pointer := my
	for len(result) < num {
		if pointer+step <= high {
			next := pointer + step
			result = append(result, "p-"+strconv.Itoa(next))
			pointer = next
		} else {
			next := low + (pointer + step - high)
			result = append(result, "p-"+strconv.Itoa(next))
			pointer = next
		}
	}

	return result
}
