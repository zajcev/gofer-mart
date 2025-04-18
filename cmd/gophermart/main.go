package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zajcev/gofer-mart/internal/gophermart/accural"
	"github.com/zajcev/gofer-mart/internal/gophermart/config"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"github.com/zajcev/gofer-mart/internal/gophermart/handlers"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	database.Migration(cfg.DatabaseURI)
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	db := database.NewDBService(pool)

	accSystem := accural.NewAccrual(db)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go accSystem.AccrualIntegration(ctx, cfg.AccSystemAddr)
	log.Fatal(http.ListenAndServe(cfg.Address, Router(db)))
}

func Router(db *database.DBService) chi.Router {
	handler := handlers.NewHandler(db)

	r := chi.NewRouter()
	//r.Use(middleware.GzipMiddleware)
	//r.Use(middleware.ZapMiddleware)

	r.Post("/api/user/register", handler.RegisterUser)
	r.Post("/api/user/login", handler.LoginUser)

	r.Get("/api/user/orders", handler.GetOrders)
	r.Post("/api/user/orders", handler.UploadOrder)

	r.Get("/api/user/balance", handler.GetBalance)
	r.Get("/api/user/withdrawals", handler.GetWithdrawals)
	r.Post("/api/user/balance/withdraw", handler.SetWithdrawals)
	return r
}
