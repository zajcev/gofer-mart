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

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) AddUser(ctx context.Context, login string, pass string) {
	m.Called(ctx, login, pass)
}

func (m *MockUserStorage) GetLogin(ctx context.Context, login string) string {
	args := m.Called(ctx, login)
	return args.String(0)
}

func (m *MockUserStorage) GetPassword(ctx context.Context, login string) string {
	args := m.Called(ctx, login)
	return args.String(0)
}

func (m *MockUserStorage) NewSession(ctx context.Context, login string, token string) error {
	args := m.Called(ctx, login, token)
	return args.Error(0)
}

func (m *MockUserStorage) GetUserID(ctx context.Context, login string) (int, error) {
	args := m.Called(ctx, login)
	return args.Int(0), args.Error(1)
}

func (m *MockUserStorage) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	args := m.Called(ctx, token)
	return args.Int(0), args.Error(1)
}

func TestRegisterUser_Success(t *testing.T) {
	mockStorage := new(MockUserStorage)
	userHandler := NewUserHandler(mockStorage)

	user := model.User{Login: "testuser", Password: "password123"}

	mockStorage.On("GetLogin", mock.Anything, user.Login).Return("")
	mockStorage.On("AddUser", mock.Anything, user.Login, mock.Anything).Return()
	mockStorage.On("NewSession", mock.Anything, user.Login, mock.Anything).Return(nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	userHandler.RegisterUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Authorization"))
	mockStorage.AssertExpectations(t)
}

func TestRegisterUser_LoginConflict(t *testing.T) {
	mockDB := new(MockUserStorage)
	userHandler := &UserHandler{db: mockDB}

	user := model.User{Login: "testuser", Password: "testpassword"}
	userJSON, _ := json.Marshal(user)

	mockDB.On("GetLogin", mock.Anything, user.Login).Return(user.Login)

	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(userJSON))
	resp := httptest.NewRecorder()

	userHandler.RegisterUser(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)
	assert.Contains(t, resp.Body.String(), "user already exists")
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "AddUser")
	mockDB.AssertNotCalled(t, "NewSession")
}

func TestLoginUser_Success(t *testing.T) {
	mockStorage := new(MockUserStorage)
	userHandler := NewUserHandler(mockStorage)

	user := model.User{Login: "testuser", Password: "password123"}
	hashedPassword, _ := hashPassword(user.Password)

	mockStorage.On("GetLogin", mock.Anything, user.Login).Return(user.Login)
	mockStorage.On("GetPassword", mock.Anything, user.Login).Return(hashedPassword)
	mockStorage.On("NewSession", mock.Anything, user.Login, mock.AnythingOfType("string")).Return(nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	userHandler.LoginUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("Authorization"))
	mockStorage.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	mockStorage := new(MockUserStorage)
	userHandler := NewUserHandler(mockStorage)

	user := model.User{Login: "testuser", Password: "wrongpassword"}
	mockStorage.On("GetLogin", mock.Anything, user.Login).Return(user.Login)
	mockStorage.On("GetPassword", mock.Anything, user.Login).Return("hashedpassword")

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	userHandler.LoginUser(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockStorage.AssertExpectations(t)
}
