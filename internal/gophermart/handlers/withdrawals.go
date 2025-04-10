package handlers

import (
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
	"net/http"
)

func GetWithdrawals(w http.ResponseWriter, r *http.Request) {
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
	list, err := database.GetWithdraw(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}
func SetWithdrawals(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserId(r.Context(), token)
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	withdraw := model.Withdraw{}
	err = json.NewDecoder(r.Body).Decode(&withdraw)
	if err != nil {
		log.Printf("Error while decode object in method SetWithdrawals: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	withdraw.UserID = userID
	resp := database.SetWithdraw(r.Context(), withdraw)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp)
}
