package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"io"
	"net/http"
	"time"
)

type OrderStorage interface {
	UploadOrder(ctx context.Context, id string, userID int, status string, uploadedAt time.Time) int
	GetOrders(ctx context.Context, userID int) ([]model.Order, error)
	GetActiveOrders(ctx context.Context) ([]model.Order, error)
}
type OrderHandler struct {
	orderStorage OrderStorage
	authStorage  AuthStorage
}

func NewOrderHandler(orderStorage OrderStorage, authStorage AuthStorage) *OrderHandler {
	return &OrderHandler{orderStorage: orderStorage, authStorage: authStorage}
}

func (oh *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	orderStorage := oh.orderStorage
	authStorage := oh.authStorage.DB

	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := authStorage.GetUserIDByToken(r.Context(), token)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	var order model.Order
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(r.Body)
	order.ID = string(body)
	if !order.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	respCode := orderStorage.UploadOrder(r.Context(), order.ID, userID, "NEW", time.Now())
	w.WriteHeader(respCode)
}

func (oh *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orderStorage := oh.orderStorage
	authStorage := oh.authStorage.DB

	token := r.Header.Get("Authorization")
	userID, err := authStorage.GetUserIDByToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	var orders []model.Order
	orders, err = orderStorage.GetOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}
