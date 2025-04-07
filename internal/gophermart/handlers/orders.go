package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func UploadOrder(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	userID := getUserIdByLogin(r.Context(), token)
	w.Header().Set("Content-Type", "application/json")
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	var order model.Order
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(r.Body)
	order.ID = string(body)
	if !order.IsValid() {
		http.Error(w, "Order format invalid", http.StatusUnprocessableEntity)
	}
	resp := database.UploadOrder(r.Context(), order.ID, userID, "", time.Now())
	w.WriteHeader(resp)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	userID := getUserIdByLogin(r.Context(), token)
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	var orders []model.Order
	orders, err := database.GetOrders(r.Context(), userID)
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

func getUserIdByLogin(ctx context.Context, login string) int {
	return database.GetUserIdByToken(ctx, login)
}
