package database

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
	"time"
)

func UploadOrder(ctx context.Context, id string, userID int, status string, uploadedAt time.Time) int {
	check := checkDuplicate(ctx, id, userID)
	switch check {
	case http.StatusOK:
		return check
	case http.StatusConflict:
		return check
	}
	return http.StatusAccepted
}

func checkDuplicate(ctx context.Context, id string, userID int) int {
	return http.StatusAccepted
}

func GetOrders(ctx context.Context, userID int) ([]model.Order, error) {
	return []model.Order{}, nil
}
