package ports

import (
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
)

type PurchaseRepository interface {
	Create(purchase *domain.Purchase) error
	GetByID(id string) (*domain.Purchase, error)
	Update(purchase *domain.Purchase) error
	Delete(id string) error
	GetAll() ([]*domain.Purchase, error)
	UpdateStatus(purchaseID string, status string, paymentID string, userID string) error
	CancelOldUnpaid(cutoff time.Time) (int64, error)
}
