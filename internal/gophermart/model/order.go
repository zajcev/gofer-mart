package model

import (
	"strconv"
	"time"
)

type Order struct {
	ID         string `json:"number"`
	UserID     int
	Status     string    `json:"status,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
	Accrual    int       `json:"accrual,omitempty"`
}

func (o *Order) IsValid() bool {
	sum := 0
	nDigits := len(o.ID)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		digit, err := strconv.Atoi(string(o.ID[i]))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
