package config

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	os.Setenv("RUN_ADDRESS", "127.0.0.1:8080")
	os.Setenv("DATABASE_URI", "postgres://testuser:testpassword@localhost:5432/testdb?sslmode=disable")
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://127.0.0.1:8090")

	config, err := NewConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.Address != "127.0.0.1:8080" {
		t.Errorf("expected Address to be '127.0.0.1:8080', got %s", config.Address)
	}
	if config.DatabaseURI != "postgres://testuser:testpassword@localhost:5432/testdb?sslmode=disable" {
		t.Errorf("expected DatabaseURI to be 'postgres://testuser:testpassword@localhost:5432/testdb?sslmode=disable', got %s", config.DatabaseURI)
	}
	if config.AccSystemAddr != "http://127.0.0.1:8090" {
		t.Errorf("expected AccSystemAddr to be 'http://127.0.0.1:8090', got %s", config.AccSystemAddr)
	}

	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("DATABASE_URI")
	os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
}
