package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
)

type MockWithdrawStorage struct {
	mock.Mock
}

func (m *MockWithdrawStorage) SetWithdraw(ctx context.Context, withdraw model.Withdraw) int {
	args := m.Called(ctx, withdraw)
	return args.Int(0)
}

func (m *MockWithdrawStorage) GetWithdraw(ctx context.Context, userID int) ([]model.Withdraw, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Withdraw), args.Error(1)
}

func (m *MockWithdrawStorage) SetBalanceWithdraw(ctx context.Context, w *model.Withdraw) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func TestGetWithdrawals_Success(t *testing.T) {
	withdrawStorage := new(MockWithdrawStorage)
	authStorage := new(MockAuthStorage)

	handler := NewWithdrawHandler(withdrawStorage, AuthStorage{DB: authStorage})

	token := "valid_token"
	userID := 1
	withdraws := []model.Withdraw{{Order: "123", UserID: userID, Sum: 100}}

	authStorage.On("GetUserIDByToken", mock.Anything, token).Return(userID, nil)
	withdrawStorage.On("GetWithdraw", mock.Anything, userID).Return(withdraws, nil)

	req, _ := http.NewRequest("GET", "/api/user/withdrawals", nil)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()

	handler.GetWithdrawals(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetWithdrawals_Unauthorized(t *testing.T) {
	withdrawStorage := new(MockWithdrawStorage)
	authStorage := new(MockAuthStorage)

	handler := NewWithdrawHandler(withdrawStorage, AuthStorage{DB: authStorage})

	req, _ := http.NewRequest("GET", "/api/user/withdrawals", nil)
	rr := httptest.NewRecorder()

	handler.GetWithdrawals(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestSetWithdrawals_Success(t *testing.T) {
	withdrawStorage := new(MockWithdrawStorage)
	authStorage := new(MockAuthStorage)

	handler := NewWithdrawHandler(withdrawStorage, AuthStorage{DB: authStorage})

	token := "valid_token"
	userID := 1
	withdraw := model.Withdraw{Order: "123", UserID: userID, Sum: 100}

	authStorage.On("GetUserIDByToken", mock.Anything, token).Return(userID, nil)
	withdrawStorage.On("SetWithdraw", mock.Anything, withdraw).Return(http.StatusCreated)
	withdrawStorage.On("SetBalanceWithdraw", mock.Anything, &withdraw).Return(nil)

	body, _ := json.Marshal(withdraw)
	req, _ := http.NewRequest("POST", "/api/user/balance/withdraw", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()

	handler.SetWithdrawals(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestSetWithdrawals_BadRequest(t *testing.T) {
	withdrawStorage := new(MockWithdrawStorage)
	authStorage := new(MockAuthStorage)

	handler := NewWithdrawHandler(withdrawStorage, AuthStorage{DB: authStorage})

	withdraw := model.Withdraw{Sum: 100}
	body, _ := json.Marshal(withdraw)
	token := "valid_token"

	authStorage.On("GetUserIDByToken", mock.Anything, token).Return(1, nil)

	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()

	handler.SetWithdrawals(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	authStorage.AssertExpectations(t)
	withdrawStorage.AssertNotCalled(t, "SetWithdraw") // Явно проверяем, что метод НЕ вызывался
}

func TestSetWithdrawals_Unauthorized(t *testing.T) {
	withdrawStorage := new(MockWithdrawStorage)
	authStorage := new(MockAuthStorage)

	handler := NewWithdrawHandler(withdrawStorage, AuthStorage{DB: authStorage})

	req, _ := http.NewRequest("POST", "/api/user/withdrawals", nil)
	rr := httptest.NewRecorder()

	handler.SetWithdrawals(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
