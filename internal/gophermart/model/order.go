package model

import (
	"github.com/spf13/cast"
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
	var sum int64
	length := 0
	n := cast.ToInt64(o.ID)
	for temp := n; temp > 0; temp /= 10 {
		length++
	}

	parity := length % 2
	pos := 0

	for n > 0 {
		digit := n % 10
		n /= 10

		if pos%2 == parity {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		pos++
	}

	return sum%10 == 0
}
