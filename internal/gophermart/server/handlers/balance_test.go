package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
)

type MockBalanceStorage struct {
	mock.Mock
}

func (m *MockBalanceStorage) GetUserBalance(ctx context.Context, userID int) (model.Balance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(model.Balance), args.Error(1)
}

func (m *MockBalanceStorage) SetCurrent(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockBalanceStorage) SetBalanceWithdraw(ctx context.Context, w *model.Withdraw) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

type MockAuthStorage struct {
	mock.Mock
}

func (m *MockAuthStorage) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	args := m.Called(ctx, token)
	return args.Int(0), args.Error(1)
}

func TestGetBalance_Success(t *testing.T) {
	mockDB := new(MockBalanceStorage)
	mockAuth := new(MockAuthStorage)
	handler := NewBalanceHandler(mockDB, AuthStorage{DB: mockAuth})

	userID := 1
	balance := model.Balance{Current: 100}
	token := "valid-token"

	mockAuth.On("GetUserIDByToken", mock.Anything, token).Return(userID, nil)
	mockDB.On("GetUserBalance", mock.Anything, userID).Return(balance, nil)

	req, err := http.NewRequest("GET", "/balance", nil)
	req.Header.Set("Authorization", token)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetBalance(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"current":100}`, rr.Body.String())
}

func TestGetBalance_Unauthorized_NoToken(t *testing.T) {
	mockDB := new(MockBalanceStorage)
	mockAuth := new(MockAuthStorage)
	handler := NewBalanceHandler(mockDB, AuthStorage{DB: mockAuth})

	req, err := http.NewRequest("GET", "/balance", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetBalance(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetBalance_Unauthorized_InvalidToken(t *testing.T) {
	mockDB := new(MockBalanceStorage)
	mockAuth := new(MockAuthStorage)
	handler := NewBalanceHandler(mockDB, AuthStorage{DB: mockAuth})

	token := "invalid-token"

	mockAuth.On("GetUserIDByToken", mock.Anything, token).Return(0, errors.New("some error"))

	req, err := http.NewRequest("GET", "/balance", nil)
	req.Header.Set("Authorization", token)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetBalance(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetBalance_InternalServerError(t *testing.T) {
	mockDB := new(MockBalanceStorage)
	mockAuth := new(MockAuthStorage)
	handler := NewBalanceHandler(mockDB, AuthStorage{DB: mockAuth})

	userID := 1
	token := "valid-token"

	mockAuth.On("GetUserIDByToken", mock.Anything, token).Return(userID, nil)
	mockDB.On("GetUserBalance", mock.Anything, userID).Return(model.Balance{}, errors.New("some error"))

	req, err := http.NewRequest("GET", "/balance", nil)
	req.Header.Set("Authorization", token)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetBalance(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
