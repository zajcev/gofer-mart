package database

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/database/scripts"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
	"time"
)

func SetWithdraw(ctx context.Context, withdraw model.Withdraw) int {
	_, err := db.Exec(ctx, scripts.SetWithdrawals, withdraw.Order, withdraw.UserID, withdraw.Sum, time.Now())
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func GetWithdraw(ctx context.Context, userID int) ([]model.Withdraw, error) {
	list := []model.Withdraw{}
	row := model.Withdraw{}
	rows, _ := db.Query(ctx, scripts.GetWithdrawals, userID)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&row.Order, &row.Sum, &row.ProcessedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, row)
	}
	return list, nil
}
