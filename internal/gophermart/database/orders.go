package database

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/database/scripts"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
	"net/http"
	"time"
)

func UploadOrder(ctx context.Context, id string, userID int, status string, uploadedAt time.Time) int {
	check := checkDuplicate(ctx, id, userID)
	if check == http.StatusAccepted {
		_, err := db.Exec(ctx, scripts.AddOrder, id, userID, status, uploadedAt, 0)
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
	rows, _ := db.Query(ctx, scripts.GetOrders, userID)
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

func GetActiveOrders(ctx context.Context) ([]model.Order, error) {
	list := []model.Order{}
	row := model.Order{}
	rows, _ := db.Query(ctx, scripts.GetActiveOrders)
	if rows.Err() != nil {
		log.Printf("Error while execute query: %v", rows.Err())
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

func UpdateOrderStatus(ctx context.Context, o *model.Order) int {
	_, err := db.Exec(ctx, scripts.UpdateOrderStatus, o.Status, o.ID)
	if err != nil {
		log.Printf("Error after updateOrderStatus: %v", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func UpdateOrderAccural(ctx context.Context, o *model.Order) int {
	_, err := db.Exec(ctx, scripts.UpdateOrderAccural, o.ID, o.Accrual)
	if err != nil {
		log.Printf("Error after updateOrderAccural: %v", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func checkDuplicate(ctx context.Context, orderID string, userID int) int {
	row, err := db.Query(ctx, scripts.GetOrder, orderID)
	if err != nil {
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
