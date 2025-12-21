package main

import (
	"cart-api/internal/config"
	"cart-api/pkg/database/postgres"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
