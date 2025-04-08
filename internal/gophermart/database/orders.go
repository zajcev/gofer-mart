package database

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
	"net/http"
	"time"
)

func UploadOrder(ctx context.Context, id string, userID int, status string, uploadedAt time.Time) int {
	check := checkDuplicate(ctx, id, userID)
	if check == http.StatusAccepted {
		_, err := db.Exec(ctx, addOrder, id, userID, status, uploadedAt, 0)
		if err != nil {
			log.Printf("Error after addOrder: %v", err)
			return http.StatusInternalServerError
		}
	}
	return check
}

func GetOrders(ctx context.Context, userID int) ([]model.Order, error) {
	list := []model.Order{}
	row := model.Order{}
	rows, _ := db.Query(ctx, getOrders, userID)
	if rows.Err() != nil {
		log.Printf("Error while execute query: %v", rows.Err())
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

func checkDuplicate(ctx context.Context, orderID string, userID int) int {
	row, err := db.Query(ctx, getOrder, orderID)
	if row.Err() != nil {
		log.Printf("Error while execute query: %v", row.Err())
		return http.StatusInternalServerError
	}
	defer row.Close()
	if !row.Next() {
		log.Printf("No order found order_id: %v", orderID)
		return http.StatusAccepted
	}
	var order model.Order
	err = row.Scan(&order.ID, &order.UserID)
	if err != nil {
		log.Printf("Error while scan token value: %v", row.Err())
		return http.StatusInternalServerError
	}
	if order.ID == orderID {
		log.Printf("Order with id %v is duplicated: %v", orderID, order.ID)
		if order.UserID != userID {
			return http.StatusConflict
		} else {
			return http.StatusOK
		}
	}
	return http.StatusAccepted
}
