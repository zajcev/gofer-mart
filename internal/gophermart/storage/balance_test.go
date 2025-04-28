package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage/scripts"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0).(pgx.Rows), ret.Error(1)
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0).(pgx.Row)
}

func (m *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0).(pgconn.CommandTag), ret.Error(1)
}

func TestSetBalanceWithdraw_Success(t *testing.T) {
	mockDB := new(MockDB)
	dbService := NewDB(mockDB)

	withdraw := &model.Withdraw{Sum: 50, UserID: 1}

	mockDB.On("Exec", mock.Anything, scripts.SetWithdraw, []interface{}{withdraw.Sum, withdraw.UserID}).
		Return(fakeUpdateCommandTag(), nil)

	err := dbService.SetBalanceWithdraw(context.Background(), withdraw)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestSetBalanceWithdraw_Error(t *testing.T) {
	mockDB := new(MockDB)
	dbService := NewDB(mockDB)

	withdraw := &model.Withdraw{Sum: 50, UserID: 1}
	expectedErr := errors.New("db error")

	mockDB.On("Exec",
		mock.Anything,
		scripts.SetWithdraw,
		[]interface{}{withdraw.Sum, withdraw.UserID},
	).Return(pgconn.CommandTag{}, expectedErr)

	err := dbService.SetBalanceWithdraw(context.Background(), withdraw)

	assert.EqualError(t, err, expectedErr.Error())
	mockDB.AssertExpectations(t)
}

func fakeUpdateCommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}
