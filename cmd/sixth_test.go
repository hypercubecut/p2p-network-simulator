package test

import (
	"fmt"
	"os"
	"p2psimulator/internal"
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/bitcoin/servicecode"
	"testing"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
)

func Test_PublishManyMinersSameBlock(t *testing.T) {
	totalInterval := time.Second * 1200

	cfg := internal.GenerateConfig(10, 15000, 0, 50,
		10, 3,
		internal.WithFullNodeBPS(int64(MBPS500)),
		internal.WithNewNodeBPS(int64(MBPS500)),
		internal.WithFullNodeDelay(100),
		internal.WithNewNodeDelay(100),
		internal.WithDifferentMinerNodeDelay("p-15001"),
		internal.WithMinerNodeBPS(int64(MBPS500)))

	simulator, err := internal.NewSimulator(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing simulator - %s", err)
		return
	}

	simulator.BuildBitcoinNetwork()

	var events []base.Event

	numberOfMiners := len(cfg.ServersCfg.AllMinersNodes)
	interval := totalInterval / time.Duration(numberOfMiners)

	for idx, miner := range cfg.ServersCfg.AllMinersNodes {
		triggerMiner := simulator.Nodes["trigger-"+miner].(*node.EndpointNode)

		mineNewBlock := triggerMiner.Send(bitcoin.NewPacket(msgtype.MineSameBlockReq,
			&bitcoin.WriteBlockRequest{BPM: 120 + idx}, nil, nil),
			simulator.SimulatorTime.Add(interval*time.Duration(idx+1)))

		events = append(events, mineNewBlock)
	}

	simulator.Run(events, "new block mining")
	simulator.Wait()

	uniqueBPM := make(map[int]int)
	uniqueLen := make(map[int]int)
	for _, nd := range simulator.Nodes {
		if n, ok := nd.(*bitcoin.Node); ok {
			if n.GetServiceCode() == servicecode.FullNode {
				bpm := n.GetChain()[len(n.GetChain())-1].BPM
				_, ok := uniqueBPM[bpm]
				if !ok {
					uniqueBPM[bpm] = 1
				} else {
					uniqueBPM[bpm] += 1
				}

				_, ok = uniqueLen[len(n.GetChain())]
				if !ok {
					uniqueLen[len(n.GetChain())] = 1
				} else {
					uniqueLen[len(n.GetChain())] += 1
				}

			}
		}
	}

	fmt.Println(fmt.Sprintf("%+v \n %d \n %+v", uniqueBPM, len(uniqueBPM), uniqueLen))
	fmt.Println(float64(uniqueBPM[120])/float64(15000)*100, "%")
}
