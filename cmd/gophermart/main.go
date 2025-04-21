package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zajcev/gofer-mart/internal/gophermart/accrual"
	"github.com/zajcev/gofer-mart/internal/gophermart/config"
	"github.com/zajcev/gofer-mart/internal/gophermart/server"
	"github.com/zajcev/gofer-mart/internal/gophermart/server/handlers"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	storage.Migration(cfg.DatabaseURI)
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURI)
	db := storage.NewDB(pool)

	handler := handlers.NewHandler(db)
	accSystem := accrual.NewAccrual(db)

	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err = accSystem.AccrualIntegration(ctx, cfg.AccSystemAddr); err != nil {
			errChan <- fmt.Errorf("accrual integration failed: %w", err)
		}
	}()

	go func() {
		log.Printf("Starting server on %s", cfg.Address)
		if err = http.ListenAndServe(cfg.Address, server.NewRouter(handler)); err != nil {
			errChan <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()

	err = <-errChan
	log.Printf("Fatal error: %v", err)
	cancel()
	os.Exit(1)
}
