package storage

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
)

func TestUpdateOrderStatus(t *testing.T) {
	dbmock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer dbmock.Close(context.Background())

	service := NewDB(dbmock)

	order := &model.Order{ID: "1", Status: "completed"}

	dbmock.ExpectExec("update orders set status = \\$1 where id = \\$2").
		WithArgs(order.Status, order.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	status := service.UpdateOrderStatus(context.Background(), order)

	if status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
