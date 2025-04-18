package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
	"net/http"
)

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserID(r.Context(), token, h)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	list, err := h.db.GetWithdraw(r.Context(), userID)
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

func (h *Handler) SetWithdrawals(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserID(r.Context(), token, h)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	withdraw := model.Withdraw{}
	err = json.NewDecoder(r.Body).Decode(&withdraw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	withdraw.UserID = userID
	resp := h.db.SetWithdraw(r.Context(), withdraw)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp)
	updateUserBalance(withdraw, h)
}

func updateUserBalance(w model.Withdraw, h *Handler) {
	err := h.db.SetBalanceWithdraw(context.Background(), &w)
	if err != nil {
		log.Printf("Error after updateUserBalance: %v", err)
	}
}
