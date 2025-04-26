package domain

import (
	"time"
)

type Purchase struct {
	ID        string
	UserID    string
	Amount    int
	CreatedAt time.Time
	Status    string
	PaymentID string
	Address   string
}
