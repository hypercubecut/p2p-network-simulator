package main

import (
	"fmt"
	"os"
	"p2psimulator/internal"
	"p2psimulator/internal/config"

	"github.com/bytedance/ns-x/v2/node"
)

func main() {
	// Todo: pass config file full path here
	cfg, err := config.NewConfigFromJson("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing config - %s", err)
		return
	}

	simulator, err := internal.NewSimulator(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing simulator - %s", err)
		return
	}

	trigger := node.NewEndpointNode()

	simulator.BuildSimpleNetwork(trigger)

	// Todo: give events here
	simulator.Run(nil)
}
