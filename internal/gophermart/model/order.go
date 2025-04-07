package model

import (
	"time"
)

type Order struct {
	ID         string    `json:"number"`
	Status     string    `json:"status,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
	Accrual    string    `json:"accrual,omitempty"`
}

func (o *Order) IsValid() bool {
	if o.ID != "" {
		return true
	}
	return false
}
