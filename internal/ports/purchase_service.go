package ports

import (
    "time"

	// "github.com/AmiraliFarazmand/PTC_Task/internal/domain"
)

type PurchaseService interface {
    CreatePurchase(userID string, amount int, address string) error
    ConfirmPayment(purchaseID string) error
    CancelUnpaidPurchases(cutoff time.Time) (int64, error)
}