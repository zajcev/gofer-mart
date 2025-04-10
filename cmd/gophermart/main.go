package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/gofer-mart/internal/gophermart/config"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"github.com/zajcev/gofer-mart/internal/gophermart/handlers"
	"log"
	"net/http"
)

func main() {
	err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = database.Init(context.Background(), config.GetDatabaseURI())
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	log.Fatal(http.ListenAndServe(config.GetAddress(), Router()))
}

func Router() chi.Router {
	r := chi.NewRouter()

	//r.Use(middleware.GzipMiddleware)
	//r.Use(middleware.ZapMiddleware)

	r.Post("/api/user/register", handlers.RegisterUser)
	r.Post("/api/user/login", handlers.LoginUser)

	r.Get("/api/user/orders", handlers.GetOrders)
	r.Post("/api/user/orders", handlers.UploadOrder)

	r.Get("/api/user/balance", handlers.GetBalance)
	r.Get("/api/user/withdrawals", handlers.GetWithdrawals)
	r.Post("/api/user/balance/withdraw", handlers.SetWithdrawals)
	return r
}
