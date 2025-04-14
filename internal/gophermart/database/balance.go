package database

import (
	"context"
	"errors"
	"github.com/zajcev/gofer-mart/internal/gophermart/database/scripts"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
)

func GetUserBalance(ctx context.Context, userID int) (model.Balance, error) {
	rowID, _ := db.Query(ctx, scripts.GetBalance, userID)
	balance := model.Balance{}
	if rowID.Err() != nil {
		return balance, rowID.Err()
	}
	defer rowID.Close()
	if !rowID.Next() {
		return balance, errors.New("not found row balance for user_id")
	}
	err := rowID.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return balance, err
	}
	return balance, nil
}

func SetCurrent(ctx context.Context, order *model.Order) error {
	_, err := db.Exec(ctx, scripts.SetBalance, order.UserID, order.Accrual)
	if err != nil {
		return err
	}
	return nil
}

func SetBalanceWithdraw(ctx context.Context, w *model.Withdraw) error {
	_, err := db.Exec(ctx, scripts.SetWithdraw, w.Sum, w.UserID)
	if err != nil {
		return err
	}
	return nil
}
