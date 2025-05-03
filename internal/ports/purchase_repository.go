package ports

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
)

type PurchaseRepository interface {
	Create(purchase *domain.Purchase) error
	Delete(id string) error
	UpdateStatus(purchaseID string, status string, paymentID string, userID string) error
	CancelPurchase(purchaseID string) error
}
