package bitcoin

import (
	"fmt"
	"p2psimulator/internal/bitcoin/servicecode"
	"sync"
	"time"

	"github.com/bytedance/ns-x/v2/base"

	"github.com/bytedance/ns-x/v2/node"
	"go.uber.org/zap"
)

const (
	defaultVersion = 7000
)

// Node is a bitcoin node implement EndpointNode
type Node struct {
	*node.EndpointNode

	// unique name in network
	name string

	serviceCode    int64
	version        int32
	seeds          []string
	availablePeers map[string]bool

	// largest local inventory
	inventory int
	// local chain
	chain []*Block

	logger *zap.Logger
	wg     sync.WaitGroup
	lc     sync.Mutex
}

func NewClientNode(name string, seeds []string, logger *zap.Logger) *Node {
	return &Node{
		EndpointNode:   node.NewEndpointNode(),
		name:           name,
		seeds:          seeds,
		logger:         logger,
		version:        defaultVersion,
		serviceCode:    servicecode.Unnamed,
		chain:          []*Block{GenesisBlock},
		availablePeers: make(map[string]bool),
		inventory:      0,
	}
}

func NewNodeWithDetails(name string, serviceCode int,
	seeds []string, logger *zap.Logger) *Node {

	switch serviceCode {
	case servicecode.NODE_NETWORK:
		return NewFullNode(name, logger)

	case servicecode.Unnamed:
		return NewClientNode(name, seeds, logger)

	default:
		return NewClientNode(name, seeds, logger)
	}
}

func NewFullNode(name string, logger *zap.Logger) *Node {
	return &Node{
		EndpointNode:   node.NewEndpointNode(),
		name:           name,
		logger:         logger,
		serviceCode:    servicecode.NODE_NETWORK,
		version:        defaultVersion,
		inventory:      MasterBlockchain[len(MasterBlockchain)-1].Index,
		chain:          MasterBlockchain,
		availablePeers: make(map[string]bool),
		wg:             sync.WaitGroup{},
		lc:             sync.Mutex{},
	}
}

func (n *Node) ID() string {
	return n.name
}

func (n *Node) GetAvailablePeers() []string {
	var val []string

	for p, _ := range n.availablePeers {
		val = append(val, p)
	}

	return val
}

func (n *Node) GetServiceCode() int64 {
	return n.serviceCode
}

func (n *Node) GetInventory() int {
	return n.inventory
}

func (n *Node) GetChain() []*Block {
	return n.chain
}

func (n *Node) AddNewPeers(peers ...string) {
	for _, p := range peers {
		if _, ok := n.availablePeers[p]; !ok {
			n.availablePeers[p] = true
		}
	}
}

// broadcast packet to all servers in the network
func (n *Node) broadcast(packet *Packet, nodes map[string]base.Node, delay time.Duration) []base.Event {
	broadcast := nodes["broadcast"].(*node.BroadcastNode)

	fmt.Println("23edx132", broadcast)

	return broadcast.Transfer(packet, time.Now().Add(delay))
}
