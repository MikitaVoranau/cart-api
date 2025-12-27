package main

import (
	"cart-api/internal/config"
	"cart-api/pkg/database/postgres"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
