package config

type serversConfig struct {
	Servers []serverDetailsConfig `json:"servers_details"`
}

type serverDetailsConfig struct {
	Name string `json:"name"`

	// output flow configs
	OutputDelayInMs int64   `json:"output_delay_in_ms"`
	OutputLossRate  float64 `json:"output_loss_rate"`

	// input flow configs
	InputDelayInMs int64   `json:"input_delay_in_ms"`
	InputLossRate  float64 `json:"input_loss_rate"`
}
