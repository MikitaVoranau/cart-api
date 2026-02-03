package app

import (
	"cart-api/internal/config"
	"cart-api/internal/repository/Cart"
	"cart-api/internal/services"
	"cart-api/internal/transport/rest"
	"cart-api/pkg/database/postgres"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Run(ctx context.Context, logger *zap.Logger) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer db.Close()

	cartRepo := Cart.New(db)
	cartService := services.NewCartService(cartRepo)
	mux := http.NewServeMux()
	cartHandler := rest.NewCartHandler(cartService, logger)

	logger.Info("starting server", zap.String("host", "localhost"), zap.String("port", "3000"))

	mux.HandleFunc("DELETE /carts/{cart_id}/items/{item_id}", cartHandler.DeleteItem)
	mux.HandleFunc("POST /carts", cartHandler.PostCart)
	mux.HandleFunc("POST /carts/{cart_id}/items", cartHandler.PostItem)
	mux.HandleFunc("GET /carts/{cart_id}", cartHandler.GetItems)
	mux.HandleFunc("GET /carts/{cart_id}/price", cartHandler.GetPrice)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: mux,
	}

	go func() {
		logger.Info("starting server", zap.String("port", cfg.HTTPPort))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

			logger.Fatal("listen and serve failed", zap.Error(err))
		}
	}()
	<-ctx.Done()

	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("server successfully shutdown...")

	return nil
}
