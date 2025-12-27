package config

import (
	"cart-api/pkg/database/postgres"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort string          `mapstructure:"HTTP_PORT"`
	Postgres postgres.Config `mapstructure:",squash"`
}

func New() (*Config, error) {
	var cfg Config
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read from config file: %w", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}
	return &cfg, nil
}
