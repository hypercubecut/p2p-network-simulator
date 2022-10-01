package bitcoin

import (
	"sync"

	"github.com/bytedance/ns-x/v2/node"
	"go.uber.org/zap"
)

// Bitcoin service code
const (
	// Unnamed This node is not a full node. It may not be able to provide any data except for the transactions it originates.
	Unnamed = 0x00

	// NODE_NETWORK This is a full node and can be asked for full blocks. It should implement all protocol
	// features available in its self-reported protocol version.
	NODE_NETWORK = 0x01

	// NODE_GETUTXO This is a full node capable of responding to the getutxo protocol request.
	// This is not supported by any currently-maintained Bitcoin node. See BIP64 for details on how this is implemented.
	NODE_GETUTXO = 0x02

	// NODE_BLOOM This is a full node capable and willing to handle bloom-filtered connections. See BIP111 for details.
	NODE_BLOOM = 0x04

	// NODE_WITNESS This is a full node that can be asked for blocks and transactions including witness data.
	// See BIP144 for details.
	NODE_WITNESS = 0x08

	// NODE_XTHIN This is a full node that supports Xtreme Thinblocks.
	// This is not supported by any currently-maintained Bitcoin node.
	NODE_XTHIN = 0x10

	// NODE_NETWORK_LIMITEDT this is the same as NODE_NETWORK but the node has at least the last 288 blocks (last 2 days).
	// See BIP159 for details on how this is implemented.
	NODE_NETWORK_LIMITED = 0x0400
)

const (
	defaultVersion = 7000
)

// Node is a bitcoin node implement EndpointNode
type Node struct {
	*node.EndpointNode

	// unique name in network
	name string

	serviceCode    uint64
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

func NewNode(name string, seeds []string, logger *zap.Logger) *Node {
	return &Node{
		EndpointNode: node.NewEndpointNode(),
		name:         name,
		seeds:        seeds,
		logger:       logger,
		version:      defaultVersion,
		serviceCode:  NODE_NETWORK,
		chain:        []*Block{genesisBlock},
		inventory:    0,
	}
}

func NewFullNode(name string, allOtherFullNodes []string, logger *zap.Logger) *Node {
	peers := map[string]bool{}
	for _, n := range allOtherFullNodes {
		peers[n] = true
	}

	return &Node{
		EndpointNode:   node.NewEndpointNode(),
		name:           name,
		logger:         logger,
		serviceCode:    NODE_NETWORK,
		version:        defaultVersion,
		availablePeers: peers,
		inventory:      0,
		chain:          nil,
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

func (n *Node) GetServiceCode() uint64 {
	return n.serviceCode
}

func (n *Node) addNewPeers(peers ...string) {
	for _, p := range peers {
		if _, ok := n.availablePeers[p]; !ok {
			n.availablePeers[p] = true
		}
	}
}
