package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfigFromJson(t *testing.T) {
	cfg, err := NewConfigFromJson("./testdata/test-config.json")
	assert.NoError(t, err)

	assert.Equal(t, true, cfg.SimulatorCfg.EnableDebugLog)
	assert.Equal(t, 3, len(cfg.ServersCfg.Servers))
	assert.Equal(t, 1, len(cfg.BitcoinCfg.ServerToSeeds))
}

func Test_NewConfigFromByte(t *testing.T) {
	payload := `
{
  "simulator": {
    "enable_debug_log": true,
    "life_time_in_min": 10
  },
  "servers": {
    "servers_details": [
      {
        "name": "p1",
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p2",
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p3",
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      }
    ]
  },
  "bitcoin": {
    "server_to_seeds": {
      "p1": ["p2", "p3"]
    }
  }
}
`
	cfg, err := NewConfigFromString(payload)
	assert.NoError(t, err)

	assert.Equal(t, true, cfg.SimulatorCfg.EnableDebugLog)
	assert.Equal(t, 3, len(cfg.ServersCfg.Servers))
	assert.Equal(t, 1, len(cfg.BitcoinCfg.ServerToSeeds))
}
