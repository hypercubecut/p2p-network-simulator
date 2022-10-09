package internal

import (
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/config"
	"testing"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/math"
	"github.com/bytedance/ns-x/v2/node"
	"github.com/bytedance/ns-x/v2/tick"
	"github.com/stretchr/testify/assert"
)

func Test_PingPongTest(t *testing.T) {
	now := time.Now()

	cfg, err := config.NewConfigFromString(simpleTestConfig)

	simulator, err := NewSimulator(cfg)
	assert.NoError(t, err)

	helper := simulator.Builder

	trigger := node.NewEndpointNode()

	network, nodes := helper.
		Chain().
		NodeWithName("restrict 1", node.NewRestrictNode(node.WithBPSLimit(1024*1024, 4*1024*1024))).
		NodeWithName("channel 1", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(150*time.Millisecond)))).
		Chain().
		NodeWithName("restrict 2", node.NewRestrictNode(node.WithPPSLimit(10, 50))).
		NodeWithName("channel 2", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(200*time.Millisecond)))).
		Chain().
		NodeWithName("peer-1", bitcoin.NewNode("192.168.0.1", []string{}, simulator.Logger)).
		Group("restrict 1", "channel 1").
		NodeWithName("peer-2", bitcoin.NewNode("192.168.0.2", []string{}, simulator.Logger)).
		Chain().
		NodeOfName("peer-2").
		Group("restrict 2", "channel 2").
		NodeOfName("peer-1").
		Chain().
		NodeWithName("trigger", trigger).
		NodeOfName("peer-1").
		Summary().
		Build()
	endpoint1 := nodes["peer-1"].(*bitcoin.Node)
	endpoint2 := nodes["peer-2"].(*bitcoin.Node)

	endpoint1.Receive(endpoint1.Handler(nodes, simulator.Logger))

	endpoint2.Receive(endpoint2.Handler(nodes, simulator.Logger))

	seedsToPeer1 := &bitcoin.Peers{
		Peers: []string{"peer-2"},
	}

	healthCheckEvent := trigger.Send(bitcoin.NewPacket(msgtype.StartMessageType, seedsToPeer1, nil, nil), now)

	network.Run([]base.Event{healthCheckEvent}, tick.NewStepClock(now, time.Second), 30*time.Second)
	defer network.Wait()
}
