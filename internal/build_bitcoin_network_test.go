package internal

import (
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"p2psimulator/internal/config"
	"testing"
	"time"

	"github.com/bytedance/ns-x/v2/node"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/stretchr/testify/assert"
)

func TestSimulator_BuildSimpleBitcoinNetWork(t *testing.T) {
	cfg, err := config.NewConfigFromString(simpleTestConfig)

	simulator, err := NewSimulator(cfg)
	assert.NoError(t, err)

	simulator.BuildBitcoinNetwork()

	peersToP1 := &bitcoin.Peers{
		Peers: []string{"p2", "p3", "p4", "p5"},
	}

	now := time.Now()

	triggerP1 := simulator.Nodes["trigger-p1"].(*node.EndpointNode)

	healthCheckEvent := triggerP1.Send(bitcoin.NewPacket(msgtype.StartMessageType, peersToP1, nil, nil), now)
	simulator.Run([]base.Event{healthCheckEvent}, "healthcheck")
	defer simulator.Network.Wait()
}

const simpleTestConfig = `
{
  "simulator": {
    "enable_debug_log": true,
    "life_time_in_min": 1
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
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p3",
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p4",
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 20000,
        "input_loss_rate": 0
      },
      {
        "name": "p5",
        "output_delay_in_ms": 200,
        "output_loss_rate": 1.0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 1.0
      }
    ]
  },
  "bitcoin": {
  }
}
`
