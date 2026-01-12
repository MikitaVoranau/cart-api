package app

import (
	"cart-api/internal/config"
	"cart-api/internal/repository/CartRepo"
	"cart-api/internal/services"
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
	cartRepo := CartRepo.New(db)
	cartService := services.NewCartService(cartRepo)
	mux := http.NewServeMux()
	fmt.Println("running")
	cartHandler := rest.NewCartHandler(cartService)
	mux.HandleFunc("DELETE /carts/{cart_id}/items/{item_id}", cartHandler.DeleteItem)
	mux.HandleFunc("POST /carts", cartHandler.PostCart)
	mux.HandleFunc("POST /carts/{cart_id}/items", cartHandler.PostItem)
	mux.HandleFunc("GET /carts/{cart_id}", cartHandler.GetItems)
	mux.HandleFunc("GET /carts/{cart_id}/price", cartHandler.GetPrice)
	if err = http.ListenAndServe("localhost:3000", mux); err != nil {
		return fmt.Errorf("run: cannot start server %w", err)
	}
	return nil
}
