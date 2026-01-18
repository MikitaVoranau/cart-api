package config

import (
	"cart-api/pkg/database/postgres"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	HTTPPort string          `mapstructure:"HTTP_PORT"`
	Postgres postgres.Config `mapstructure:",squash"`
}

func New() (*Config, error) {
	var cfg Config
	viper.AutomaticEnv()

	_ = viper.BindEnv("HTTP_PORT")
	_ = viper.BindEnv("POSTGRES_HOST")
	_ = viper.BindEnv("POSTGRES_PORT")
	_ = viper.BindEnv("POSTGRES_USER")
	_ = viper.BindEnv("POSTGRES_PASS")
	_ = viper.BindEnv("POSTGRES_DB")

	viper.SetConfigFile(".env")

	if _, err := os.Stat(".env"); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("cannot read from config file: %w", err)
		}
	} else {
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}

	return &cfg, nil
}
