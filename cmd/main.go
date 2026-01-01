package main

import (
	"cart-api/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("cannot create service %w", err)
	}
}
