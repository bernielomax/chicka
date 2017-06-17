package exec

import (
	"github.com/spf13/viper"
)

type Config struct {
	Checks Checks `json:"checks"`
}

func ReadConfig() (*Config, error) {

	cfg := Config{}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) Refresh() error {

	new := Config{}

	err := viper.Unmarshal(&new)
	if err != nil {
		return err
	}

	*cfg = new

	return nil
}
