package bitcoin

import (
	"p2psimulator/internal/bitcoin/nodestate"
	"sync"

	"github.com/bytedance/ns-x/v2/node"
	"go.uber.org/zap"
)

// Node implement EndpointNode
type Node struct {
	*node.EndpointNode

	name   string // unique name in network
	logger *zap.Logger

	// bitcoin related property
	State    nodestate.State
	DNSSeeds []string

	// full nodes
	availablePeers []string

	wg sync.WaitGroup
	lc sync.Mutex
}

func NewNode(name string, seeds []string, logger *zap.Logger) *Node {
	return &Node{
		EndpointNode: node.NewEndpointNode(),
		name:         name,
		DNSSeeds:     seeds,
		logger:       logger,
		State:        nodestate.New,
	}
}

func (n *Node) ID() string {
	return n.name
}
