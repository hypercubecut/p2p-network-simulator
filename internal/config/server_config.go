package config

import (
	"encoding/json"
	"errors"
)

type ServersConfig struct {
	Servers        []*ServerDetailsConfig `json:"servers_details"`
	AllFullNodes   []string               `json:"all_full_nodes"`
	AllNewNodes    []string               `json:"all_new_nodes"`
	AllMinersNodes []string               `json:"all_miners_nodes"`
}

type ServerDetailsConfig struct {
	Name        string `json:"name"`
	ServiceCode int64  `json:"service_code"`
	Version     int    `json:"version"`

	// for client bitcoin node only
	Seeds []string `json:"seeds"`

	Peers []string `json:"peers"`

	// output flow configs
	OutputDelayInMs int64   `json:"output_delay_in_ms"`
	OutputLossRate  float64 `json:"output_loss_rate"`
	OutputBPS       int64   `json:"output_bps"`
	OutputPPS       int64   `json:"output_pps"`

	// input flow configs
	InputDelayInMs int64   `json:"input_delay_in_ms"`
	InputLossRate  float64 `json:"input_loss_rate"`
	InputBPS       int64   `json:"input_bps"`
	InputPPS       int64   `json:"input_pps"`

	QueueLimit int64 `json:"queue_limit"`
}

func (s *ServerDetailsConfig) UnmarshalJSON(data []byte) error {
	res := &struct {
		Name            string   `json:"name"`
		ServiceCode     int64    `json:"service_code"`
		Version         int64    `json:"version"`
		Seeds           []string `json:"seeds"`
		OutputDelayInMs int64    `json:"output_delay_in_ms"`
		OutputLossRate  float64  `json:"output_loss_rate"`
		OutputBPS       int64    `json:"output_bps"`
		OutputPPS       int64    `json:"output_pps"`
		InputDelayInMs  int64    `json:"input_delay_in_ms"`
		InputLossRate   float64  `json:"input_loss_rate"`
		InputBPS        int64    `json:"input_bps"`
		InputPPS        int64    `json:"input_pps"`
		QueueLimit      int64    `json:"queue_limit"`
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

	if res.InputBPS == 0 {
		res.InputBPS = 1024 * 1024
	}

	if res.InputPPS == 0 {
		res.InputPPS = 50
	}

	if res.OutputBPS == 0 {
		res.OutputBPS = 1024 * 1024
	}

	if res.OutputPPS == 0 {
		res.OutputPPS = 50
	}

	if res.QueueLimit == 0 {
		res.QueueLimit = 10
	}

	s.Name = res.Name
	s.ServiceCode = res.ServiceCode
	s.Version = int(res.Version)
	s.Seeds = res.Seeds
	s.OutputDelayInMs = res.OutputDelayInMs
	s.OutputLossRate = res.OutputLossRate
	s.OutputBPS = res.OutputBPS
	s.OutputPPS = res.OutputPPS
	s.InputDelayInMs = res.InputDelayInMs
	s.InputLossRate = res.InputLossRate
	s.InputBPS = res.InputBPS
	s.InputPPS = res.InputPPS
	s.QueueLimit = res.QueueLimit

	return nil
}
