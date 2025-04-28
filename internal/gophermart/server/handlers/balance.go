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
	balanceStorage BalanceStorage
	authStorage    AuthStorage
}

func NewBalanceHandler(balanceStorage BalanceStorage, authStorage AuthStorage) *BalanceHandler {
	return &BalanceHandler{balanceStorage: balanceStorage, authStorage: authStorage}
}

func (bh *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	authStorage := bh.authStorage.DB
	balanceStorage := bh.balanceStorage

	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := authStorage.GetUserIDByToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := balanceStorage.GetUserBalance(r.Context(), userID)
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
