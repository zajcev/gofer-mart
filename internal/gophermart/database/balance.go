package database

import (
	"context"
	"errors"
	"github.com/zajcev/gofer-mart/internal/gophermart/database/scripts"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
)

func GetUserBalance(ctx context.Context, userID int) (model.Balance, error) {
	rowId, _ := db.Query(ctx, scripts.GetBalance, userID)
	balance := model.Balance{}
	if rowId.Err() != nil {
		log.Printf("Error while execute query: %v", rowId.Err())
		return balance, rowId.Err()
	}
	defer rowId.Close()
	if !rowId.Next() {
		log.Printf("Not found row balance for user_id: %v", userID)
		return balance, errors.New("not found row balance for user_id")
	}
	err := rowId.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		log.Printf("Error while scan balance object: %v", rowId.Err())
		return balance, err
	}
	return balance, nil
}

func SetCurrent(ctx context.Context, order *model.Order) error {
	_, err := db.Exec(ctx, scripts.SetBalance, order.UserID, order.Accrual)
	if err != nil {
		log.Printf("Error after setCurrent: %v", err)
		return err
	}
	return nil
}

func SetBalanceWithdraw(ctx context.Context, w *model.Withdraw) error {
	_, err := db.Exec(ctx, scripts.SetWithdraw, w.Sum, w.UserID)
	if err != nil {
		log.Printf("Error after SetBalanceWithdraw DB: %v", err)
		return err
	}
	return nil
}
