package app

import (
	"cart-api/internal/config"
	"cart-api/internal/repository/CartRepo"
	"cart-api/internal/services"
	"cart-api/internal/transport/rest"
	"cart-api/pkg/database/postgres"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

// добавить логгер, запустить через докер, сделать тестирование

func Run(logger *zap.Logger) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer db.Close()
	cartRepo := CartRepo.New(db)
	cartService := services.NewCartService(cartRepo)
	mux := http.NewServeMux()
	cartHandler := rest.NewCartHandler(cartService, logger)
	logger.Info("starting server", zap.String("host", "localhost"), zap.String("port", "3000"))
	mux.HandleFunc("DELETE /carts/{cart_id}/items/{item_id}", cartHandler.DeleteItem)
	mux.HandleFunc("POST /carts", cartHandler.PostCart)
	mux.HandleFunc("POST /carts/{cart_id}/items", cartHandler.PostItem)
	mux.HandleFunc("GET /carts/{cart_id}", cartHandler.GetItems)
	mux.HandleFunc("GET /carts/{cart_id}/price", cartHandler.GetPrice)
	if err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), mux); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	return nil
}
