package handlers

import (
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"log"
	"net/http"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserId(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := database.GetUserBalance(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting balance: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	log.Println(balance)
	resp, err := json.Marshal(&balance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}
