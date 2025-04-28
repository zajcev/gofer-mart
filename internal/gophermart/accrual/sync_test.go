package accrual

import (
	"context"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"io"
	"net/http"
	"testing"
	"time"
)

type mockAccrualStorage struct {
	activeOrders []model.Order
}

func (m *mockAccrualStorage) UpdateOrderStatus(ctx context.Context, o *model.Order) int {
	return 1
}

func (m *mockAccrualStorage) UpdateOrderAccrual(ctx context.Context, o *model.Order) int {
	return 1
}

func (m *mockAccrualStorage) GetActiveOrders(ctx context.Context) ([]model.Order, error) {
	return m.activeOrders, nil
}

func (m *mockAccrualStorage) SetCurrent(ctx context.Context, order *model.Order) error {
	return nil
}

func TestAccrualIntegration_Success(t *testing.T) {
	mockDB := &mockAccrualStorage{
		activeOrders: []model.Order{{ID: "1"}},
	}
	acc := NewAccrual(mockDB)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go acc.AccrualIntegration(ctx, "http://localhost:8080")

	time.Sleep(3 * time.Second)

	if len(mockDB.activeOrders) != 1 {
		t.Errorf("expected 1 active order but got %d", len(mockDB.activeOrders))
	}
}

func TestSendToAccrualSystem_Success(t *testing.T) {
	order := &model.Order{ID: "1"}
	url := "http://localhost:8080"

	mockResponse := `{"ID": "1", "Status": "PROCESSED"}`
	http.HandleFunc("/api/orders/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, mockResponse)
	})

	go http.ListenAndServe(":8080", nil)

	result, err := sendToAccrualSystem(order, url)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if result.ID != order.ID {
		t.Errorf("expected order ID %s but got %s", order.ID, result.ID)
	}
}

func TestSendToAccrualSystem_ErrorResponse(t *testing.T) {
	order := &model.Order{ID: "2"}
	url := "http://localhost:8080"

	http.HandleFunc("/api/orders/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	go http.ListenAndServe(":8080", nil)

	_, err := sendToAccrualSystem(order, url)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestSendToAccrualSystem_EmptyResponse(t *testing.T) {
	order := &model.Order{ID: "3"}
	url := "http://localhost:8080"

	http.HandleFunc("/api/orders/3", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "")
	})

	go http.ListenAndServe(":8080", nil)

	_, err := sendToAccrualSystem(order, url)
	if err == nil {
		t.Fatalf("expected empty response error but got %v", err)
	}
}
