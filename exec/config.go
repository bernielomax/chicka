package exec

import (
	"github.com/spf13/viper"
)

// Config is the data structure for chick base configurations.
type Config struct {
	Plugins Plugins `json:"plugins"`
	Logging Logging `json:"logging"`
	Checks  Checks  `json:"checks"`
}

// Plugins is the data structure for configuring pluigns.
type Plugins struct {
	Path string `json:"path"`
}

type Logging struct {
	Path string `json:"path"`
}

// ReadConfig reads the chicka config file into a data structure.
func ReadConfig() (*Config, error) {

	cfg := Config{}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Refresh reloads the config data structure from the chicka config file.
func (cfg *Config) Refresh() error {

	new := Config{}

	err := viper.Unmarshal(&new)
	if err != nil {
		return err
	}

	*cfg = new

	return nil
}
