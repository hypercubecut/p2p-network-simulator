package main

import (
	"fmt"
	"os"
	"p2psimulator/internal"
	"p2psimulator/internal/config"
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

	simulator.BuildBitcoinNetwork()

	// Todo: give events here
	simulator.Run(nil)
}
