package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
	"net/http"
)

type WithdrawStorage interface {
	SetWithdraw(ctx context.Context, withdraw model.Withdraw) int
	GetWithdraw(ctx context.Context, userID int) ([]model.Withdraw, error)
	SetBalanceWithdraw(ctx context.Context, w *model.Withdraw) error
}

type WithdrawHandler struct {
	db   WithdrawStorage
	auth *UserHandler
}

func NewWithdrawHandler(db WithdrawStorage, auth *UserHandler) *WithdrawHandler {
	return &WithdrawHandler{db: db, auth: auth}
}

func (wh *WithdrawHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := wh.auth.getUserID(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	list, err := wh.db.GetWithdraw(r.Context(), userID)
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

func (wh *WithdrawHandler) SetWithdrawals(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := wh.auth.getUserID(r.Context(), token)
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
	resp := wh.db.SetWithdraw(r.Context(), withdraw)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp)
	updateUserBalance(withdraw, wh)
}

func updateUserBalance(w model.Withdraw, h *WithdrawHandler) {
	err := h.db.SetBalanceWithdraw(context.Background(), &w)
	if err != nil {
		log.Printf("Error after updateUserBalance: %v", err)
	}
}
