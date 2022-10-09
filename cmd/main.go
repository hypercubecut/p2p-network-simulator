package main

import (
	"fmt"
	"math/rand"
	"os"
	"p2psimulator/internal"
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
)

func main() {
	cfg := internal.GenerateConfig(1000, 200, 10, 1)

	simulator, err := internal.NewSimulator(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing simulator - %s", err)
		return
	}

	simulator.BuildBitcoinNetwork()

	var events []base.Event

	for _, peer := range cfg.ServersCfg.AllNewNodes {
		triggerPeer := simulator.Nodes["trigger-"+peer].(*node.EndpointNode)

		i := rand.Intn(10)

		event := triggerPeer.Send(
			bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil),
			time.Now().Add(time.Second*time.Duration(i)))

		events = append(events, event)
	}

	simulator.Run(events, "peer discovery")
	simulator.Wait()

	triggerMiner := simulator.Nodes["trigger-m1"].(*node.EndpointNode)

	events = []base.Event{}
	i := 1
	for i < 21 {
		mineNewBlock := triggerMiner.Send(bitcoin.NewPacket(msgtype.MineNewBlockReq,
			&bitcoin.WriteBlockRequest{BPM: 120}, nil, nil),
			time.Now().Add(time.Minute*time.Duration(i*5)))

		events = append(events, mineNewBlock)

		i++
	}

	simulator.Run(events, "new block mining")
	simulator.Wait()
}
