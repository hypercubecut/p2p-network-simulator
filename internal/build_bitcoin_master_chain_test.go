package internal

import (
	"p2psimulator/internal/bitcoin"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulator_BuildBitcoinMasterChain(t *testing.T) {
	buildInitialMasterChain(10)

	assert.Equal(t, 10, len(bitcoin.MasterBlockchain))
}
