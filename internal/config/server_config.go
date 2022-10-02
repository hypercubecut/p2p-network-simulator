package config

import (
	"encoding/json"
	"errors"
)

type serversConfig struct {
	Servers []serverDetailsConfig `json:"servers_details"`
}

type serverDetailsConfig struct {
	Name        string `json:"name"`
	ServiceCode int64  `json:"service_code"`
	Version     int    `json:"version"`

	// for client bitcoin node only
	Seeds []string `json:"seeds"`

	// output flow configs
	OutputDelayInMs int64   `json:"output_delay_in_ms"`
	OutputLossRate  float64 `json:"output_loss_rate"`

	// input flow configs
	InputDelayInMs int64   `json:"input_delay_in_ms"`
	InputLossRate  float64 `json:"input_loss_rate"`
}

func (s *serverDetailsConfig) UnmarshalJSON(data []byte) error {
	res := &struct {
		Name            string   `json:"name"`
		ServiceCode     int64    `json:"service_code"`
		Version         int64    `json:"version"`
		Seeds           []string `json:"seeds"`
		OutputDelayInMs int64    `json:"output_delay_in_ms"`
		OutputLossRate  float64  `json:"output_loss_rate"`
		InputDelayInMs  int64    `json:"input_delay_in_ms"`
		InputLossRate   float64  `json:"input_loss_rate"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	if res.Name == "" {
		return errors.New("name cannot be empty")
	}

	if res.Version == 0 {
		res.Version = 7000
	}

	if res.InputDelayInMs == 0 {
		res.InputDelayInMs = 200
	}

	if res.OutputDelayInMs == 0 {
		res.OutputDelayInMs = 200
	}

	s.Name = res.Name
	s.ServiceCode = res.ServiceCode
	s.Version = int(res.Version)
	s.Seeds = res.Seeds
	s.OutputDelayInMs = res.OutputDelayInMs
	s.OutputLossRate = res.OutputLossRate
	s.InputDelayInMs = res.InputDelayInMs
	s.InputLossRate = res.InputLossRate

	return nil
}
