package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"io"
	"net/http"
	"time"
)

func (h *Handler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserID(r.Context(), token, h)
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
	respCode := h.db.UploadOrder(r.Context(), order.ID, userID, "NEW", time.Now())
	w.WriteHeader(respCode)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	userID, err := getUserID(r.Context(), token, h)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	var orders []model.Order
	orders, err = h.db.GetOrders(r.Context(), userID)
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

func getUserID(ctx context.Context, token string, h *Handler) (int, error) {
	userID, err := h.db.GetUserIDByToken(ctx, token)
	return userID, err
}
