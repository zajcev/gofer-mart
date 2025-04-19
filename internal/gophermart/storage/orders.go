package storage

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage/scripts"
	"net/http"
	"time"
)

func (s *DBService) UploadOrder(ctx context.Context, id string, userID int, status string, uploadedAt time.Time) int {
	check := checkDuplicate(ctx, id, userID, s)
	if check == http.StatusAccepted {
		_, err := s.db.Exec(ctx, scripts.AddOrder, id, userID, status, uploadedAt, 0)
		if err != nil {
			return http.StatusInternalServerError
		}
	}
	return check
}

func (s *DBService) GetOrders(ctx context.Context, userID int) ([]model.Order, error) {
	var list []model.Order
	row := model.Order{}
	rows, _ := s.db.Query(ctx, scripts.GetOrders, userID)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&row.ID, &row.Status, &row.UploadedAt, &row.Accrual)
		if err != nil {
			return nil, err
		}
		list = append(list, row)
	}
	return list, nil
}

func (s *DBService) GetActiveOrders(ctx context.Context) ([]model.Order, error) {
	list := []model.Order{}
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

func (s *DBService) UpdateOrderStatus(ctx context.Context, o *model.Order) int {
	_, err := s.db.Exec(ctx, scripts.UpdateOrderStatus, o.Status, o.ID)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *DBService) UpdateOrderAccural(ctx context.Context, o *model.Order) int {
	_, err := s.db.Exec(ctx, scripts.UpdateOrderAccural, o.ID, o.Accrual)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func checkDuplicate(ctx context.Context, orderID string, userID int, s *DBService) int {
	row, err := s.db.Query(ctx, scripts.GetOrder, orderID)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer row.Close()
	if !row.Next() {
		return http.StatusAccepted
	}
	var order model.Order
	err = row.Scan(&order.ID, &order.UserID)
	if err != nil {
		return http.StatusInternalServerError
	}
	if order.ID == orderID {
		if order.UserID != userID {
			return http.StatusConflict
		} else {
			return http.StatusOK
		}
	}
	return http.StatusAccepted
}
