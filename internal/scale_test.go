package internal

import (
	"fmt"
	"math/rand"
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"testing"
	time "time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"

	"github.com/stretchr/testify/assert"
)

func TestSimulator_Scale(t *testing.T) {
	cfg := GenerateConfig(100, 2000, 20, 0)

	sim, err := NewSimulator(cfg)
	assert.NoError(t, err)

	sim.BuildBitcoinNetwork()

	var events []base.Event

	for _, peer := range cfg.ServersCfg.AllNewNodes {
		triggerPeer := sim.Nodes["trigger-"+peer].(*node.EndpointNode)

		i := rand.Intn(10)

		event := triggerPeer.Send(
			bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil),
			time.Now().Add(time.Second*time.Duration(i)))

		events = append(events, event)
	}

	sim.Run(events, "peer discovery")
	sim.Wait()

	for _, peer := range cfg.ServersCfg.AllNewNodes {
		n := sim.Nodes[peer].(*bitcoin.Node)
		assert.Equal(t, 100, len(n.GetChain()), fmt.Sprintf("node %s", peer))
	}
}
