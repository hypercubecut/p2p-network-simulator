package config

type SimulatorConfig struct {
	EnableDebugLog bool  `json:"enable_debug_log"`
	LifeTimeInMin  int64 `json:"life_time_in_min"`
}
