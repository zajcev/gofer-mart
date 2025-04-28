package storage

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage/scripts"
)

type MockAuthDB struct {
	mock.Mock
}

func (m *MockAuthDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	args = append([]interface{}{ctx, sql}, args...)
	returnArgs := m.Called(args...)
	return returnArgs.Get(0).(pgx.Rows), returnArgs.Error(1)
}

func (m *MockAuthDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	args = append([]interface{}{ctx, sql}, args...)
	returnArgs := m.Called(args...)
	return returnArgs.Get(0).(pgconn.CommandTag), returnArgs.Error(1)
}

func (m *MockAuthDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	args = append([]interface{}{ctx, sql}, args...)
	returnArgs := m.Called(args...)
	return returnArgs.Get(0).(pgx.Row)
}

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Close() {
	m.Called()
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockRows) Columns() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func TestAddUser(t *testing.T) {
	mockDB := new(MockAuthDB)
	service := NewDB(mockDB)

	ctx := context.Background()
	login := "testuser"
	pass := "password"

	mockDB.On("Exec", ctx, scripts.AddUser, login, pass).Return(pgconn.CommandTag{}, nil).Once()

	service.AddUser(ctx, login, pass)

	mockDB.AssertExpectations(t)
}
