package model

import "time"

type Withdraw struct {
	Order       string    `json:"order"`
	UserID      int       `json:"user_id,omitempty"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
