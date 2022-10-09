package internal

import (
	"p2psimulator/internal/bitcoin"
	"time"

	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/math"
	"github.com/bytedance/ns-x/v2/node"
)

var (
	dropPacketCallback = func(packet base.Packet, source, target base.Node, now time.Time) {
		logger.Info("emit packet")
	}
)

func (s *Simulator) BuildBitcoinNetwork() {
	buildInitialMasterChain(s.cfg.BitcoinCfg.MasterChainLen)

	router := node.NewScatterNode(node.WithRouteSelector(func(packet base.Packet, nodes []base.Node) base.Node {
		if p, ok := packet.(*bitcoin.Packet); ok {
			return s.Nodes[genRestrictInName(p.Destination.ID())]
		}
		panic("no route to host")
	}))

	broadcast := node.NewBroadcastNode()

	var allServers []*bitcoin.Node
	for _, serverCfg := range s.cfg.ServersCfg.Servers {
		server := bitcoin.NewNodeWithDetails(serverCfg.Name,
			int(serverCfg.ServiceCode), s.cfg.ServersCfg.AllFullNodes, s.Logger)

		// add trigger node for each bitcoin node
		s.Builder.Chain().
			NodeWithName("trigger-"+serverCfg.Name, node.NewEndpointNode()).
			NodeWithName(server.ID(), server)

		// ChannelNode is a simulated network channel with loss, delay and reorder features
		var outChannelOpt []node.Option
		if serverCfg.OutputDelayInMs != 0 {
			outChannelOpt = append(outChannelOpt, node.WithDelay(math.NewFixedDelay(time.Millisecond*time.Duration(serverCfg.OutputDelayInMs))))
		}
		if serverCfg.OutputLossRate != 0 {
			outChannelOpt = append(outChannelOpt, node.WithTransferCallback(dropPacketCallback),
				node.WithLoss(math.NewRandomLoss(serverCfg.OutputLossRate, s.Random)))
		}
		channelNodeOut := node.NewChannelNode(outChannelOpt...)

		// RestrictNode simulate a node with limited ability
		// Once packets through a RestrictNode reaches the limit(in bps or pps), the later packets will be put in a queue
		// Once the queue overflow, later packets will be discarded
		restrictNodeOut := node.NewRestrictNode(
			node.WithBPSLimit(float64(serverCfg.OutputBPS), serverCfg.QueueLimit*serverCfg.OutputBPS),
			node.WithPPSLimit(float64(serverCfg.OutputPPS), serverCfg.QueueLimit*serverCfg.OutputPPS))

		// output flow chain
		// server -> channel -> restrict -> router
		s.Builder.Chain().NodeWithName(server.ID(), server).
			NodeWithName(genChannelOutName(server.ID()), channelNodeOut).
			NodeWithName(genRestrictOutName(server.ID()), restrictNodeOut).
			NodeWithName("router", router)

		// input flow chain
		var inChannelOpt []node.Option
		if serverCfg.InputDelayInMs != 0 {
			inChannelOpt = append(inChannelOpt, node.WithDelay(math.NewFixedDelay(time.Millisecond*time.Duration(serverCfg.InputDelayInMs))))
		}
		if serverCfg.InputLossRate != 0 {
			inChannelOpt = append(inChannelOpt, node.WithTransferCallback(dropPacketCallback),
				node.WithLoss(math.NewRandomLoss(serverCfg.InputLossRate, s.Random)))
		}
		channelNodeIn := node.NewChannelNode(inChannelOpt...)

		restrictNodeIn := node.NewRestrictNode(
			node.WithBPSLimit(float64(serverCfg.InputBPS), serverCfg.QueueLimit*serverCfg.InputBPS),
			node.WithPPSLimit(float64(serverCfg.InputPPS), serverCfg.QueueLimit*serverCfg.InputPPS))

		// broadcast -> restrict
		s.Builder.Chain().NodeWithName("broadcast", broadcast).
			NodeWithName(genRestrictInName(server.ID()), restrictNodeIn)

		// router -> restrict -> channel -> server
		s.Builder.Chain().NodeWithName("router", router).
			NodeWithName(genRestrictInName(server.ID()), restrictNodeIn).
			NodeWithName(genChannelInName(server.ID()), channelNodeIn).
			NodeWithName(server.ID(), server)

		allServers = append(allServers, server)
	}

	network, nodes := s.Builder.Build()

	for _, server := range allServers {
		server.Receive(server.Handler(nodes, s.Logger))
	}

	s.Network = network
	s.Nodes = nodes
}

func buildInitialMasterChain(masterChainLen int) {
	var i int

	for i < masterChainLen-1 {
		newBlock, err := bitcoin.GenerateBlock(bitcoin.MasterBlockchain[len(bitcoin.MasterBlockchain)-1], 55)
		if err != nil {
			panic("failed to generate initial blocks")
		}

		// Skip bitcoin.IsBlockValid() here

		bitcoin.MasterBlockchain = append(bitcoin.MasterBlockchain, newBlock)

		i++
	}
}
