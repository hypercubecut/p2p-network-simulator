package internal

import (
	"math/rand"
	"os"
	"p2psimulator/internal/config"
	"sync"
	"time"

	nsx "github.com/bytedance/ns-x/v2"
	"github.com/bytedance/ns-x/v2/base"
	"github.com/bytedance/ns-x/v2/tick"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	logger     *zap.Logger
	onceLogger sync.Once
)

type Simulator struct {
	Logger *zap.Logger

	cfg *config.Config

	LifeTime time.Duration

	Topology *simple.WeightedUndirectedGraph

	Builder nsx.Builder
	Random  *rand.Rand
	Network *nsx.Network
	Nodes   map[string]base.Node
}

func NewSimulator(cfg *config.Config) (*Simulator, error) {
	file, err := os.Create("log.txt")
	if err != nil {
		return nil, err
	}

	logInitOnce(cfg.SimulatorCfg.EnableDebugLog, file)

	source := rand.NewSource(0)
	random := rand.New(source)
	helper := nsx.NewBuilder()

	return &Simulator{
		Logger:   logger,
		LifeTime: time.Duration(cfg.SimulatorCfg.LifeTimeInMin) * time.Minute,
		cfg:      cfg,
		Builder:  helper,
		Random:   random,
	}, nil
}

func (s *Simulator) Run(events []base.Event) {
	s.Network.Run(events, tick.NewStepClock(time.Now(), time.Second), s.LifeTime)
}

func genChannelOutName(id string) string {
	return "channel-out-" + id
}

func genChannelInName(id string) string {
	return "channel-in-" + id
}

func genRestrictOutName(id string) string {
	return "restrict-out-" + id
}

func genRestrictInName(id string) string {
	return "restrict-in-" + id
}
