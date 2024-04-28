package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Environments define the environment variables
type Environments struct {
	ApiPort         string   `mapstructure:"PORT" envconfig:"PORT"`
	ApiToken        string   `mapstructure:"API_TOKEN" envconfig:"API_TOKEN"`
	ApiBaseUrls     []string `mapstructure:"API_BASE_URLS" envconfig:"API_BASE_URLS"`
	PollingInterval int      `mapstructure:"POLLING_INTERVAL" envconfig:"POLLING_INTERVAL"`
}

// LoadEnvVars load the environment variables
func LoadEnvVars() (*Environments, error) {
	godotenv.Load()
	c := &Environments{}
	if err := envconfig.Process("", c); err != nil {
		return nil, err
	}
	return c, nil
}
