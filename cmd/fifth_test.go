package test

import (
	"fmt"
	"os"
	"p2psimulator/internal"
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"testing"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
)

func Test_P2PBroadcast(t *testing.T) {
	cfg := internal.GenerateConfig(10, 15000, 0, 0,
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

	triggerNode := simulator.Nodes["trigger-p-15000"].(*node.EndpointNode)

	event := triggerNode.Send(bitcoin.NewPacket(msgtype.P2PWithBroadcastMessageType,
		&bitcoin.P2PMessage{
			MsgID: "34523243525",
			Dest:  "p-10000"}, nil, nil),
		simulator.SimulatorTime)

	simulator.Run(base.Aggregate(event), "P2P Broadcasting")
	simulator.Wait()

}
