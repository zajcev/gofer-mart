package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/gofer-mart/internal/gophermart/middleware"
	"github.com/zajcev/gofer-mart/internal/gophermart/server/handlers"
)

func NewRouter(handler *handlers.Handler) chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.ZapMiddleware)

	r.Post("/api/user/register", handler.RegisterUser)
	r.Post("/api/user/login", handler.LoginUser)

	r.Get("/api/user/orders", handler.GetOrders)
	r.Post("/api/user/orders", handler.UploadOrder)

	r.Get("/api/user/balance", handler.GetBalance)
	r.Get("/api/user/withdrawals", handler.GetWithdrawals)
	r.Post("/api/user/balance/withdraw", handler.SetWithdrawals)
	return r
}
