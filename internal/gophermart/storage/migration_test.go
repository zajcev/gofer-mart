package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	mockDB := &MockDB{}
	dbService := NewDB(mockDB)

	assert.NotNil(t, dbService)
	assert.Equal(t, mockDB, dbService.db)
}
