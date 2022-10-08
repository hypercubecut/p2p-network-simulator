package internal

import (
	"fmt"
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/config"
	"testing"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"github.com/stretchr/testify/assert"
)

func TestSimulator_TestNodesDiscovery(t *testing.T) {
	cfg, err := config.NewConfigFromString(threeBitcoinClientNodesTestConfig)

	fmt.Printf("%+v \n", cfg)

	simulator, err := NewSimulator(cfg)
	assert.NoError(t, err)

	simulator.BuildBitcoinNetwork()

	now := time.Now()

	triggerP1 := simulator.Nodes["trigger-p1"].(*node.EndpointNode)

	triggerP1PeerDiscovery :=
		triggerP1.Send(bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil), now)

	triggerP2 := simulator.Nodes["trigger-p2"].(*node.EndpointNode)
	triggerP2PeerDiscovery :=
		triggerP2.Send(bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil), now)

	triggerP3 := simulator.Nodes["trigger-p3"].(*node.EndpointNode)
	triggerP3PeerDiscovery :=
		triggerP3.Send(bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil), now)

	simulator.Run([]base.Event{triggerP1PeerDiscovery, triggerP2PeerDiscovery, triggerP3PeerDiscovery}, "peer discovery")
	simulator.Network.Wait()

	p1 := simulator.Nodes["p1"].(*bitcoin.Node)
	fmt.Println("p1's available peers", p1.GetAvailablePeers())

	p2 := simulator.Nodes["p2"].(*bitcoin.Node)
	fmt.Println("p2's available peers", p2.GetAvailablePeers())

	p3 := simulator.Nodes["p3"].(*bitcoin.Node)
	fmt.Println("p3's available peers", p3.GetAvailablePeers())
}

const threeBitcoinClientNodesTestConfig = `
{
  "simulator": {
    "enable_debug_log": true,
    "life_time_in_min": 2
  },
  "servers": {
    "servers_details": [
      {
        "name": "p1",
		"seeds": ["p2", "p3"],
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p2",
		"seeds": ["p1", "p3"],
        "output_delay_in_ms": 100,
        "output_loss_rate": 0,
        "input_delay_in_ms": 100,
        "input_loss_rate": 0
      },
	  {
        "name": "p3",
        "output_delay_in_ms": 500,
        "output_loss_rate": 0,
        "input_delay_in_ms": 500,
        "input_loss_rate": 0
      }
    ]
  },
  "bitcoin": {
  }
}
`
