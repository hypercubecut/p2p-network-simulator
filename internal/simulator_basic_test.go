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

func Test_BasicSanityTest(t *testing.T) {
	cfg, err := config.NewConfigFromString(simpleTestConfig)

	simulator, err := NewSimulator(cfg)
	assert.NoError(t, err)

	helper := simulator.Builder
	now := time.Now()
	routeTable := make(map[base.Node]base.Node)
	ipTable := make(map[string]base.Node)
	scatter := node.NewScatterNode(node.WithRouteSelector(func(packet base.Packet, nodes []base.Node) base.Node {
		if p, ok := packet.(*bitcoin.Packet); ok {
			return routeTable[p.Destination]
		}
		panic("no route to host")
	}))

	client := bitcoin.NewClientNode("127.0.0.1", []string{}, simulator.Logger)

	network, nodes := helper.
		Chain().
		NodeWithName("client", client).
		NodeWithName("scatter", scatter).
		NodeWithName("route1", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(time.Millisecond*200)))).
		NodeWithName("server1", bitcoin.NewClientNode("192.168.0.1", []string{}, simulator.Logger)).
		Chain().
		Node(client).
		Node(scatter).
		NodeWithName("route2", node.NewChannelNode(node.WithDelay(math.NewFixedDelay(time.Millisecond*300)))).
		NodeWithName("server2", bitcoin.NewClientNode("192.168.0.2", []string{}, simulator.Logger)).Summary().
		Build()

	server1 := nodes["server1"].(*bitcoin.Node)
	server2 := nodes["server2"].(*bitcoin.Node)
	route1 := nodes["route1"]
	route2 := nodes["route2"]
	routeTable[server1] = route1
	routeTable[server2] = route2
	ipTable["192.168.0.1"] = server1
	ipTable["192.168.0.2"] = server2
	server1.Receive(server1.Handler(nodes, simulator.Logger)) // server 1 should receive after 1-second send delay + 200 milliseconds channel delay
	server2.Receive(server2.Handler(nodes, simulator.Logger)) // server 2 should receive after 2-second send delay + 200 milliseconds channel delay
	sender := createSender(client, ipTable)
	events := make([]base.Event, 0)
	events = append(events, sender(&bitcoin.Packet{}, "192.168.0.1", now.Add(time.Second*1))) // send to server1 after 1 second
	events = append(events, sender(&bitcoin.Packet{}, "192.168.0.2", now.Add(time.Second*2))) // send to server2 after 2 second
	network.Run(events, tick.NewStepClock(now, time.Millisecond), 300*time.Second)
	defer network.Wait()
}

type sender func(packet *bitcoin.Packet, ip string, t time.Time) base.Event

func createSender(client *bitcoin.Node, ipTable map[string]base.Node) sender {
	return func(packet *bitcoin.Packet, ip string, t time.Time) base.Event {
		return client.Send(bitcoin.NewPacket(msgtype.StartMessageType, packet.GetPayload(), client, ipTable[ip].(*bitcoin.Node)), t)
	}
}
