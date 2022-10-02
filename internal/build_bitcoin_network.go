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
	router := node.NewScatterNode(node.WithRouteSelector(func(packet base.Packet, nodes []base.Node) base.Node {
		if p, ok := packet.(*bitcoin.Packet); ok {
			return s.Nodes[genRestrictInName(p.Destination.ID())]
		}
		panic("no route to host")
	}))

	var allServers []*bitcoin.Node
	for _, serverCfg := range s.cfg.ServersCfg.Servers {
		server := bitcoin.NewNodeWithDetails(serverCfg.Name,
			int(serverCfg.ServiceCode), serverCfg.Seeds, s.Logger)

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
		restrictNodeOut := node.NewRestrictNode(node.WithBPSLimit(1024*1024, 4*1024*1024))

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

		restrictNodeIn := node.NewRestrictNode(node.WithBPSLimit(1024*1024, 4*1024*1024))

		// router -> restrict -> channel -> server
		s.Builder.Chain().NodeWithName("router", router).
			NodeWithName(genRestrictInName(server.ID()), restrictNodeIn).
			NodeWithName(genChannelInName(server.ID()), channelNodeIn).
			NodeWithName(server.ID(), server)

		allServers = append(allServers, server)
	}

	network, nodes := s.Builder.Summary().Build()

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
