package storage

import (
	"context"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
	"testing"
	"time"
)

func TestUpload_Failed(t *testing.T) {
	dbmock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer dbmock.Close(context.Background())

	order := &model.Order{ID: "1234567890123456", UserID: 1, Status: "NEW", UploadedAt: time.Now()}
	service := NewDB(dbmock)
	result := service.UploadOrder(context.Background(), order.ID, order.UserID, order.Status, order.UploadedAt)

	assert.Equal(t, http.StatusInternalServerError, result)
}
