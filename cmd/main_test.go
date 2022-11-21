package test

import (
	"fmt"
	"os"
	"p2psimulator/internal"
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/bitcoin/servicecode"
	"testing"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
)

var (
	MBPS500 = 62500000
)

func Test_PublishOneNewBlock(t *testing.T) {
	cfg := internal.GenerateConfig(10, 15000, 0, 1,
		1, 2,
		internal.WithFullNodeBPS(int64(MBPS500)),
		internal.WithNewNodeBPS(int64(MBPS500)),
		internal.WithFullNodeDelay(100),
		internal.WithNewNodeDelay(100),
		internal.WithMinerNodeDelay(100),
		internal.WithMinerNodeBPS(int64(MBPS500)))

	simulator, err := internal.NewSimulator(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing simulator - %s", err)
		return
	}

	simulator.BuildBitcoinNetwork()

	triggerMiner := simulator.Nodes["trigger-p-15001"].(*node.EndpointNode)

	var events []base.Event
	mineNewBlock := triggerMiner.Send(bitcoin.NewPacket(msgtype.MineNewBlockReq,
		&bitcoin.WriteBlockRequest{BPM: 120}, nil, nil), simulator.SimulatorTime)

	events = append(events, mineNewBlock)

	simulator.Run(events, "new block mining")
	simulator.Wait()

	c := 0
	for _, nd := range simulator.Nodes {
		if n, ok := nd.(*bitcoin.Node); ok {
			if n.GetServiceCode() == servicecode.FullNode {
				if len(n.GetChain()) == 11 {
					c += 1
				}
			}
		}
	}

	fmt.Println(float64(c) / 15000.0 * 100)
}
