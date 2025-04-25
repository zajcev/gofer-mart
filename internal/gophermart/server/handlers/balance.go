package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
)

type BalanceStorage interface {
	GetUserBalance(ctx context.Context, userID int) (model.Balance, error)
	SetCurrent(ctx context.Context, order *model.Order) error
	SetBalanceWithdraw(ctx context.Context, w *model.Withdraw) error
}

type BalanceHandler struct {
	db   BalanceStorage
	auth *AuthStorage
}

func NewBalanceHandler(db BalanceStorage, auth *AuthStorage) *BalanceHandler {
	return &BalanceHandler{db: db, auth: auth}
}

func (bh *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := bh.auth.db.GetUserIDByToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := bh.db.GetUserBalance(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
