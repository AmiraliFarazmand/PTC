package ports

import (
	"time"
)

type PurchaseService interface {
	CreatePurchase(userID string, amount int, address string) (string, error)
	ConfirmPayment(purchaseID string, userID string) error
	CancelUnpaidPurchases(cutoff time.Time) (int64, error)
}
