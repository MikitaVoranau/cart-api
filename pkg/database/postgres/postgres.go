package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST" mapstructure:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT" mapstructure:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER" mapstructure:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASS" mapstructure:"POSTGRES_PASS"`
	Database string `env:"POSTGRES_DB" mapstructure:"POSTGRES_DB"`
}

func New(cfg *Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("cant ping to postgres: %w", err)
	}
	return conn, nil
}
