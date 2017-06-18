package exec

import (
	"time"

	"github.com/spf13/viper"
)

// Config is the data structure for chick base configurations.
type Config struct {
	Plugins PluginsConfig `json:"plugins"`
	Logging LoggingConfig `json:"logging"`
	Cache   CacheConfig   `json:"cache"`
	Checks  Checks        `json:"checks"`
}

// PluginsConfig is the data structure for configuring plugins.
type PluginsConfig struct {
	Path string `json:"path"`
}

// LoggingConfig is the data structure for configuring logging.
type LoggingConfig struct {
	Path string `json:"path"`
}

// CacheConfig is the data structure for configuring the go-cache.
type CacheConfig struct {
	TTL time.Duration `json:"ttl"`
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
