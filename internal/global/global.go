package global

import (
	"sync"

	"yourapp/pkg/config"
)

var (
	configOnce sync.Once
	appConfig  *config.Config
)

// Config represents the application configuration
type Config struct {
	App      config.AppConfig      `yaml:"app"`
	Server   config.ServerConfig   `yaml:"server"`
	Database config.DatabaseConfig `yaml:"database"`
	Cache    config.CacheConfig    `yaml:"cache"`
	Logging  config.LoggingConfig  `yaml:"logging"`
}

// Init initializes the global configuration
func Init() {
	// This function can be used for any global initialization
	// Currently, configuration is loaded in bootstrap.Start()
}

// SetConfig sets the global configuration
func SetConfig(cfg *config.Config) {
	configOnce.Do(func() {
		appConfig = cfg
	})
}

// GetConfig returns the global configuration
func GetConfig() *config.Config {
	return appConfig
}
