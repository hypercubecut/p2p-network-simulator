package test

import (
	"fmt"
	"math/rand"
	"os"
	"p2psimulator/internal"
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"testing"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
)

func Test_SecondMain(t *testing.T) {
	cfg := internal.GenerateConfig(40000, 10, 1, 0,
		10, 2,
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

	var events []base.Event

	for _, peer := range cfg.ServersCfg.AllNewNodes {
		fmt.Println(len(cfg.ServersCfg.Servers), peer)
		triggerPeer := simulator.Nodes["trigger-"+peer].(*node.EndpointNode)

		i := rand.Intn(10)

		event := triggerPeer.Send(
			bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil),
			time.Now().Add(time.Second*time.Duration(i)))

		events = append(events, event)
	}

	simulator.Run(events, "peer discovery")
	simulator.Wait()

	//triggerMiner := simulator.Nodes["trigger-p-15001"].(*node.EndpointNode)
	//
	//events = []base.Event{}
	//mineNewBlock := triggerMiner.Send(bitcoin.NewPacket(msgtype.MineNewBlockReq,
	//	&bitcoin.WriteBlockRequest{BPM: 120}, nil, nil), time.Now())
	//
	//events = append(events, mineNewBlock)
	//
	//simulator.Run(events, "new block mining")
	//simulator.Wait()
	//
	//c := 0
	//for _, nd := range simulator.Nodes {
	//	if n, ok := nd.(*bitcoin.Node); ok {
	//		if n.GetServiceCode() == servicecode.FullNode {
	//			if len(n.GetChain()) == 11 {
	//				c += 1
	//			}
	//		}
	//	}
	//}
	//
	//fmt.Println(float64(c) / 15000.0 * 100)
}
