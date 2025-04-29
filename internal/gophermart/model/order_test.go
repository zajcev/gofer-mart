package model

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		orderID  string
		expected bool
	}{
		{"123456789", false},
		{"123456788", false},
		{"1234567a89", false},
		{"000000000", true},
		{"111111111", false},
		{"123", false},
	}

	for _, test := range tests {
		order := Order{ID: test.orderID}
		result := order.IsValid()
		if result != test.expected {
			t.Errorf("IsValid() for ID %s = %v; expected %v", test.orderID, result, test.expected)
		}
	}
}
