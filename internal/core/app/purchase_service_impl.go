package app

import (
	"fmt"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
)

type PurchaseServiceImpl struct {
	PurchaseRepo ports.PurchaseRepository
}

var _ ports.PurchaseService = (*PurchaseServiceImpl)(nil) // Verify interface implementation

func (s *PurchaseServiceImpl) CreatePurchase(userID string, amount int, address string) (string, error) {
	if amount <= 0 || address == "" {
		return "", fmt.Errorf("invalid purchase details")
	}

	purchaseID := utils.GenerateRandomID()

	purchase := &domain.Purchase{
		ID:        purchaseID,
		UserID:    userID,
		Amount:    amount,
		CreatedAt: time.Now(),
		Status:    "pending",
		PaymentID: "",
		Address:   address,
	}
	if err := s.PurchaseRepo.Create(purchase); err != nil {
		return "", err
	}

	return purchaseID, nil
}

func (s *PurchaseServiceImpl) ConfirmPayment(purchaseID string, userID string) error {
	paymentID := generatePaymentID(purchaseID)
	return s.PurchaseRepo.UpdateStatus(purchaseID, "packaging and delivering", paymentID, userID)
}

func (s *PurchaseServiceImpl) CancelUnpaidPurchase(purchaseID string) error {
	return s.PurchaseRepo.CancelPurchase(purchaseID)
}

func generatePaymentID(purchaseID string) string {
	return "PAY-" + purchaseID + utils.RandSeq(10)
}
