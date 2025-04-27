package ports

import (
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
)

type PurchaseRepository interface {
	Create(purchase *domain.Purchase) error
	Delete(id string) error
	UpdateStatus(purchaseID string, status string, paymentID string, userID string) error
	CancelOldUnpaid(cutoff time.Time) (int64, error)
}
