package database

import (
	"context"
	"errors"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
)

func GetUserBalance(ctx context.Context, userID int) (model.Balance, error) {
	rowId, _ := db.Query(ctx, getBalance, userID)
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
	err := rowId.Scan(&balance)
	if err != nil {
		log.Printf("Error while scan balance object: %v", rowId.Err())
		return balance, err
	}
	return balance, nil
}
