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

func TestSimulator_TestBlocksFirst(t *testing.T) {
	cfg, err := config.NewConfigFromString(twoBitcoinNodesTestConfig)

	fmt.Printf("%+v \n", cfg)

	buildInitialMasterChain(cfg.BitcoinCfg.MasterChainLen)

	simulator, err := NewSimulator(cfg)
	assert.NoError(t, err)

	simulator.BuildBitcoinNetwork()

	now := time.Now()

	triggerP1 := simulator.Nodes["trigger-p1"].(*node.EndpointNode)

	triggerP1PeerDiscovery :=
		triggerP1.Send(bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil), now)

	simulator.Run([]base.Event{triggerP1PeerDiscovery})
	simulator.Network.Wait()

	p1 := simulator.Nodes["p1"].(*bitcoin.Node)
	//fmt.Println("p1's inventory is", p1.GetInventory())
	//spew.Dump(p1.GetChain())

	p2 := simulator.Nodes["p2"].(*bitcoin.Node)
	//fmt.Println("p2's inventory is", p2.GetInventory())

	assert.Equal(t, p2.GetInventory(), p1.GetInventory())
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
		"seeds": ["p2"],
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p2",
		"service_code": 1,
        "output_delay_in_ms": 100,
        "output_loss_rate": 0,
        "input_delay_in_ms": 100,
        "input_loss_rate": 0
      }
    ]
  },
  "bitcoin": {
	"master_chain_len": 1000
  }
}
`
