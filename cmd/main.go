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
	cfg := internal.GenerateConfig(1000, 200, 10, 0)

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

	//events = []base.Event{}
	//
	//for _, peer := range cfg.ServersCfg.AllNewNodes {
	//	i := rand.Intn(10)
	//
	//	triggerPeer := simulator.Nodes["trigger-"+peer].(*node.EndpointNode)
	//	mineNewBlockEvent := triggerPeer.Send(bitcoin.NewPacket(msgtype.MineNewBlockReq,
	//		&bitcoin.WriteBlockRequest{BPM: i * 13}, nil, nil), time.Now().
	//		Add(time.Second*time.Duration(i*5)))
	//
	//	events = append(events, mineNewBlockEvent)
	//}
	//
	//simulator.Run(events, "new block mining")
	//simulator.Wait()
}
