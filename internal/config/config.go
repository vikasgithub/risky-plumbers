package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/vikasgithub/risky-plumbers/internal/log"
	"os"
)

const (
	defaultServerPort = 8080
)

// Config represents an application configuration.
type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port int
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(file string, logger log.Logger) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	configBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	viper.SetDefault("server.port", defaultServerPort)
	var c Config
	if err := viper.MergeConfig(bytes.NewBuffer(configBytes)); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		logger.Fatal("Environment can't be loaded: ", err)
	}

	return &c, nil
}
