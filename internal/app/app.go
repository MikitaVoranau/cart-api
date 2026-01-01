package app

import (
	"cart-api/internal/config"
	"cart-api/internal/transport/rest"
	"cart-api/pkg/database/postgres"
	"fmt"
	"net/http"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("run: cannot connect config %w", err)
	}
	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		return fmt.Errorf("run: cannot connect to database %w", err)
	}
	defer db.Close()
	router := rest.NewRouter()

	if err = http.ListenAndServe("localhost:3000", router); err != nil {
		return fmt.Errorf("run: cannot start server %w", err)
	}
	return nil
}
