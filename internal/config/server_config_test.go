package config

import (
	"encoding/json"
	"p2psimulator/internal/bitcoin/servicecode"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerDetailsConfig_UnmarshalJSON(t *testing.T) {
	testJson := []byte(`{
"servers_details": [
      {
        "name": "p1",
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      },
      {
        "name": "p2"
      },
      {
        "name": "p3",
		"service_code": 1,
        "output_delay_in_ms": 200,
        "output_loss_rate": 0,
        "input_delay_in_ms": 200,
        "input_loss_rate": 0
      }
    ]
}`)

	cfg := &ServersConfig{}

	err := json.Unmarshal(testJson, cfg)
	assert.NoError(t, err)

	assert.Equal(t, "p2", cfg.Servers[1].Name)
	assert.Equal(t, 200, int(cfg.Servers[2].InputDelayInMs))

	assert.Equal(t, "p3", cfg.Servers[2].Name)
	assert.Equal(t, servicecode.FullNode, int(cfg.Servers[2].ServiceCode))
}
