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
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserId(r.Context(), token)
	log.Printf("User ID: %v, Token: %v\n", userID, token)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	var order model.Order
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(r.Body)
	order.ID = string(body)
	if !order.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(database.UploadOrder(r.Context(), order.ID, userID, "NEW", time.Now()))
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	userID, err := getUserId(r.Context(), token)
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	var orders []model.Order
	orders, err = database.GetOrders(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting orders: %v", err)
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

func getUserId(ctx context.Context, token string) (int, error) {
	userID, err := database.GetUserIdByToken(ctx, token)
	return userID, err
}
