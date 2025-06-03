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
	withdrawStorage WithdrawStorage
	authStorage     AuthStorage
}

func NewWithdrawHandler(withdrawStorage WithdrawStorage, authStorage AuthStorage) *WithdrawHandler {
	return &WithdrawHandler{withdrawStorage: withdrawStorage, authStorage: authStorage}
}

func (wh *WithdrawHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	withdrawStorage := wh.withdrawStorage
	authStorage := wh.authStorage.DB

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
	list, err := withdrawStorage.GetWithdraw(r.Context(), userID)
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
	withdrawStorage := wh.withdrawStorage
	authStorage := wh.authStorage.DB

	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := authStorage.GetUserIDByToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userID == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	withdraw := model.Withdraw{}
	err = json.NewDecoder(r.Body).Decode(&withdraw)
	if withdraw.Order == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	withdraw.UserID = userID
	resp := withdrawStorage.SetWithdraw(r.Context(), withdraw)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp)
	updateUserBalance(withdraw, wh)
}

func updateUserBalance(w model.Withdraw, h *WithdrawHandler) {
	err := h.withdrawStorage.SetBalanceWithdraw(context.Background(), &w)
	if err != nil {
		log.Printf("Error after updateUserBalance: %v", err)
	}
}
