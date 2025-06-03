package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
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
	order := &model.Order{ID: "1", Status: "PROCESSING"}
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

func TestUpdateOrderAccrual(t *testing.T) {
	dbmock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer dbmock.Close(context.Background())

	service := NewDB(dbmock)
	order := &model.Order{ID: "1", Accrual: 100}
	dbmock.ExpectExec("update orders set accural = \\$2 where id = \\$1").
		WithArgs(order.ID, order.Accrual).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	status := service.UpdateOrderAccrual(context.Background(), order)

	if status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSetCurrent(t *testing.T) {
	dbmock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbmock.Close(context.Background())

	order := &model.Order{ID: "1", Accrual: 100}
	expectedSQL := `INSERT INTO balance .*`
	dbmock.ExpectExec(expectedSQL).
		WithArgs(order.UserID, order.Accrual).
		WillReturnResult(pgxmock.NewResult("INSERT", 1)) // 1 строка обновлена
	service := NewDB(dbmock)
	err = service.SetCurrent(context.Background(), order)
	assert.NoError(t, err)
	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
