package storage_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage"
	"testing"
)

type mockRows struct {
	called bool
	value  interface{}
}

func (r *mockRows) Close()     {}
func (r *mockRows) Err() error { return nil }

func (r *mockRows) Next() bool {
	if !r.called {
		r.called = true
		return true
	}
	return false
}

func (r *mockRows) Scan(dest ...interface{}) error {
	if len(dest) == 0 {
		return errors.New("no destination arguments provided")
	}
	d := dest[0]
	switch v := d.(type) {
	case *int:
		if val, ok := r.value.(int); ok {
			*v = val
		}
	case *string:
		if val, ok := r.value.(string); ok {
			*v = val
		}
	default:
		return fmt.Errorf("unsupported destination type: %T", d)
	}

	return nil
}

func (r *mockRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 1")
}

func (r *mockRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{}
}

func (r *mockRows) RawValues() [][]byte {
	return [][]byte{}
}

func (r *mockRows) Values() ([]interface{}, error) {
	return []interface{}{r.value}, nil
}

func (r *mockRows) Conn() *pgx.Conn {
	var conn pgx.Conn
	return &conn
}

type mockDB struct {
	login string
}

func (m *mockDB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	if args[0] == m.login {
		return &mockRows{value: m.login}, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("INSERT 1"), nil
}

func (m *mockDB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return &mockRows{value: 42}
}

func (m *mockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func (m *mockDB) Commit(ctx context.Context) error {
	return nil
}

func (m *mockDB) Rollback(ctx context.Context) error {
	return nil
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	db := &mockDB{login: "testuser"}
	s := storage.NewDB(db)
	s.AddUser(ctx, "testuser", "test")
}

func TestNewSession(t *testing.T) {
	ctx := context.Background()
	db := &mockDB{login: "testuser"}
	s := storage.NewDB(db)
	err := s.NewSession(ctx, "testuser", "sometoken")
	require.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	db := &mockDB{login: "testuser"}
	s := storage.NewDB(db)
	s.GetLogin(ctx, "testuser")
}

func TestGetUserID(t *testing.T) {
	ctx := context.Background()
	db := &mockDB{login: "testuser"}
	s := storage.NewDB(db)
	s.GetUserID(ctx, "testuser")
}
