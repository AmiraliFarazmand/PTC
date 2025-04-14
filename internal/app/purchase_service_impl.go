package app

import (
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/domain"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
)

type PurchaseServiceImpl struct {
	PurchaseRepo domain.PurchaseRepository
}

func (s *PurchaseServiceImpl) CreatePurchase(userID string, amount int, address string) (string, error) {

	purchaseID := utils.GenerateRandomID()

	purchase := domain.Purchase{
		ID:        purchaseID,
		UserID:    userID,
		Amount:    amount,
		CreatedAt: time.Now(),
		Status:    "pending",
		PaymentID: "",
		Address:   address,
	}

	return purchaseID, s.PurchaseRepo.Create(purchase)
}

func (s *PurchaseServiceImpl) ConfirmPayment(purchaseID string, userID string) error {
	paymentID := generatePaymentID(purchaseID)
	return s.PurchaseRepo.UpdateStatus(purchaseID, "packaging and delivering", paymentID, userID)
}

func (s *PurchaseServiceImpl) CancelUnpaidPurchases(cutoff time.Time) (int64, error) {
	return s.PurchaseRepo.CancelOldUnpaid(cutoff)
}

func generatePaymentID(purchaseID string) string {
	return "PAY-" + purchaseID + utils.RandSeq(10)
}
