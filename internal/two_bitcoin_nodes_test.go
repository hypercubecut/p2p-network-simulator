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

func TestSimulator_TestTwoNodesDiscovery(t *testing.T) {
	cfg, err := config.NewConfigFromString(twoBitcoinNodesTestConfig)

	simulator, err := NewSimulator(cfg)
	assert.NoError(t, err)

	trigger := node.NewEndpointNode()

	simulator.BuildSimpleNetwork(trigger)

	now := time.Now()
	triggerP1PeerDiscovery :=
		trigger.Send(bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil), now)

	simulator.Run([]base.Event{triggerP1PeerDiscovery})
	simulator.Network.Wait()

	p1 := simulator.Nodes["p1"].(*bitcoin.Node)
	fmt.Println("p1's available peers", p1.GetAvailablePeers())

	p2 := simulator.Nodes["p2"].(*bitcoin.Node)
	fmt.Println("p2's available peers", p2.GetAvailablePeers())

	p3 := simulator.Nodes["p3"].(*bitcoin.Node)
	fmt.Println("p3's available peers", p3.GetAvailablePeers())
}

const twoBitcoinNodesTestConfig = `
{
  "simulator": {
    "enable_debug_log": true,
    "life_time_in_min": 2
  },
  "servers": {
    "servers_details": [
      {
        "name": "p1",
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p2",
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
    "server_to_seeds": {
      "p1": ["p2", "p3"]
    }
  }
}
`
