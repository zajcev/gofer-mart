package storage

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage/scripts"
	"net/http"
)

func (s *DBService) SetCurrent(ctx context.Context, order *model.Order) error {
	_, err := s.db.Exec(ctx, scripts.SetBalance, order.UserID, order.Accrual)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBService) GetActiveOrders(ctx context.Context) ([]model.Order, error) {
	var list []model.Order
	row := model.Order{}
	rows, _ := s.db.Query(ctx, scripts.GetActiveOrders)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&row.ID, &row.UserID, &row.Status, &row.Accrual)
		if err != nil {
			return nil, err
		}
		list = append(list, row)
	}
	return list, nil
}

func (s *DBService) UpdateOrderAccrual(ctx context.Context, o *model.Order) int {
	_, err := s.db.Exec(ctx, scripts.UpdateOrderAccural, o.ID, o.Accrual)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *DBService) UpdateOrderStatus(ctx context.Context, o *model.Order) int {
	_, err := s.db.Exec(ctx, scripts.UpdateOrderStatus, o.Status, o.ID)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
