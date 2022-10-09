package config

type Config struct {
	loader *Loader

	SimulatorCfg *SimulatorConfig
	ServersCfg   *ServersConfig

	// bitcoin related config
	BitcoinCfg *BitcoinConfig
}

func NewConfigFromJson(file string) (*Config, error) {
	loader, err := newJsonConfig(file)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		loader: loader,
	}

	cfg.loadConfig()

	return cfg, nil
}

func NewConfigFromString(payload string) (*Config, error) {
	loader, err := newFromBytes([]byte(payload))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		loader: loader,
	}

	cfg.loadConfig()

	return cfg, nil
}

func (c *Config) loadConfig() {
	err := c.loader.GetJSON("simulator", &c.SimulatorCfg)
	if err != nil {
		// setup default value
		c.SimulatorCfg = &SimulatorConfig{
			EnableDebugLog: false,
			LifeTimeInMin:  10,
		}
	}

	err = c.loader.GetJSON("servers", &c.ServersCfg)
	if err != nil || len(c.ServersCfg.Servers) == 0 {
		panic("failed to load server config err:" + err.Error())
	}

	_ = c.loader.GetJSON("bitcoin", &c.BitcoinCfg)
}
