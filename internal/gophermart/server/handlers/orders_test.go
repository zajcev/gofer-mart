package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
)

type MockOrderStorage struct {
	mock.Mock
}

func (m *MockOrderStorage) UploadOrder(ctx context.Context, id string, userID int, status string, uploadedAt time.Time) int {
	args := m.Called(ctx, id, userID, status, uploadedAt)
	return args.Int(0)
}

func (m *MockOrderStorage) GetOrders(ctx context.Context, userID int) ([]model.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderStorage) GetActiveOrders(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Order), args.Error(1)
}

func TestUploadOrder_Success(t *testing.T) {
	orderStorage := new(MockOrderStorage)
	authStorage := new(MockAuthStorage)
	orderHandler := NewOrderHandler(orderStorage, AuthStorage{DB: authStorage})

	authStorage.On("GetUserIDByToken", mock.Anything, "valid_token").Return(1, nil)
	orderStorage.On("UploadOrder", mock.Anything, "403204", 1, "NEW", mock.Anything).Return(http.StatusOK)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(`403204`))
	req.Header.Set("Authorization", "valid_token")
	w := httptest.NewRecorder()

	orderHandler.UploadOrder(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	orderStorage.AssertCalled(t, "UploadOrder", mock.Anything, "403204", 1, "NEW", mock.Anything)
}

func TestUploadOrder_Unauthorized(t *testing.T) {
	orderStorage := new(MockOrderStorage)
	authStorage := new(MockAuthStorage)
	orderHandler := NewOrderHandler(orderStorage, AuthStorage{DB: authStorage})

	req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBufferString(`{"id":"12351325135"}`))
	w := httptest.NewRecorder()

	orderHandler.UploadOrder(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetOrders_Success(t *testing.T) {
	orderStorage := new(MockOrderStorage)
	authStorage := new(MockAuthStorage)
	orderHandler := NewOrderHandler(orderStorage, AuthStorage{DB: authStorage})

	authStorage.On("GetUserIDByToken", mock.Anything, "valid_token").Return(1, nil)
	orderStorage.On("GetOrders", mock.Anything, 1).Return([]model.Order{{ID: "order_id"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.Header.Set("Authorization", "valid_token")
	w := httptest.NewRecorder()

	orderHandler.GetOrders(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var orders []model.Order
	err := json.Unmarshal(w.Body.Bytes(), &orders)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, "order_id", orders[0].ID)
}

func TestGetOrders_Unauthorized(t *testing.T) {
	orderStorage := new(MockOrderStorage)
	authStorage := new(MockAuthStorage)
	authStorage.On("GetUserIDByToken", mock.Anything, "").Return(0, errors.New("unauthorized"))

	orderHandler := NewOrderHandler(orderStorage, AuthStorage{DB: authStorage})

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	orderHandler.GetOrders(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	authStorage.AssertExpectations(t)
}
