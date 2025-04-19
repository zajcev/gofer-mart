package storage

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage/scripts"
	"net/http"
	"time"
)

func (s *DBService) SetWithdraw(ctx context.Context, withdraw model.Withdraw) int {
	_, err := s.db.Exec(ctx, scripts.SetWithdrawals, withdraw.Order, withdraw.UserID, withdraw.Sum, time.Now())
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *DBService) GetWithdraw(ctx context.Context, userID int) ([]model.Withdraw, error) {
	list := []model.Withdraw{}
	row := model.Withdraw{}
	rows, _ := s.db.Query(ctx, scripts.GetWithdrawals, userID)
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
