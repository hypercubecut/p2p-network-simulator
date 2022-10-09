package internal

import (
	"p2psimulator/internal/bitcoin"
	"p2psimulator/internal/bitcoin/msgtype"
	"testing"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/node"
	"github.com/stretchr/testify/assert"
)

func TestSimulator_TestBlocksFirst(t *testing.T) {
	cfg := GenerateConfig(1000, 1, 1, 1)

	simulator, err := NewSimulator(cfg)
	assert.NoError(t, err)

	simulator.BuildBitcoinNetwork()

	now := time.Now()

	triggerP1 := simulator.Nodes["trigger-p1"].(*node.EndpointNode)

	triggerP1PeerDiscovery :=
		triggerP1.Send(bitcoin.NewPacket(msgtype.PeerDiscoveryMessageType, nil, nil, nil), now)

	simulator.Run([]base.Event{triggerP1PeerDiscovery}, "p1 peer discovery")
	simulator.Wait()

	p1 := simulator.Nodes["p1"].(*bitcoin.Node)
	//fmt.Println("p1's inventory is", p1.GetInventory())
	//spew.Dump(p1.GetChain())

	p2 := simulator.Nodes["f1"].(*bitcoin.Node)
	//fmt.Println("p2's inventory is", p2.GetInventory())
	//spew.Dump(p2.GetChain())

	assert.Equal(t, p2.GetInventory(), p1.GetInventory())

	triggerP3 := simulator.Nodes["trigger-m1"].(*node.EndpointNode)

	var events []base.Event
	i := 1
	for i < 21 {
		p3MineNewBlock := triggerP3.Send(bitcoin.NewPacket(msgtype.MineNewBlockReq,
			&bitcoin.WriteBlockRequest{BPM: 120}, nil, nil),
			now.Add(time.Second*time.Duration(i*5)))

		events = append(events, p3MineNewBlock)

		i++
	}

	simulator.Run(events, "new block mining")
	simulator.Wait()

	assert.Equal(t, p2.GetInventory(), p1.GetInventory())

	assert.Equal(t, 1020, len(p1.GetChain()))
	assert.Equal(t, 1020, len(p2.GetChain()))
}
