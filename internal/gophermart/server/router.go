package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/gofer-mart/internal/gophermart/middleware"
	"github.com/zajcev/gofer-mart/internal/gophermart/server/handlers"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage"
)

func NewRouter(db *storage.DBService) chi.Router {

	auth := handlers.NewUserHandler(db)
	order := handlers.NewOrderHandler(db, auth)
	balance := handlers.NewBalanceHandler(db, auth)
	withdraw := handlers.NewWithdrawHandler(db, auth)

	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.ZapMiddleware)

	r.Post("/api/user/register", auth.RegisterUser)
	r.Post("/api/user/login", auth.LoginUser)

	r.Get("/api/user/orders", order.GetOrders)
	r.Post("/api/user/orders", order.UploadOrder)

	r.Get("/api/user/balance", balance.GetBalance)
	r.Get("/api/user/withdrawals", withdraw.GetWithdrawals)
	r.Post("/api/user/balance/withdraw", withdraw.SetWithdrawals)
	return r
}
