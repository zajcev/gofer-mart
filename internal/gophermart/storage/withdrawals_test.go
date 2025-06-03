package storage

import (
	"context"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
	"testing"
)

func TestUploadWithdraw_Failed(t *testing.T) {
	dbmock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbmock.Close(context.Background())

	withdraw := model.Withdraw{Order: "1234567890123456", UserID: 1, Sum: 100.0}
	service := NewDB(dbmock)
	result := service.SetWithdraw(context.Background(), withdraw)

	assert.Equal(t, http.StatusInternalServerError, result)
}
