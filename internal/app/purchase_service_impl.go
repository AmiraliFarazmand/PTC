package app

import (
    "errors"
    "time"

    "github.com/AmiraliFarazmand/PTC_Task/internal/domain"
    // "github.com/AmiraliFarazmand/PTC_Task/internal/ports"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type PurchaseServiceImpl struct {
    PurchaseRepo domain.PurchaseRepository
}

func (s *PurchaseServiceImpl) CreatePurchase(userID string, amount int, address string) (string, error) {
    userObjectID, err := bson.ObjectIDFromHex(userID)
    if err != nil {
        return "", errors.New("invalid user ID format")
    }

    purchase := domain.Purchase{
        ID:        bson.NewObjectID(),
        UserID:    userObjectID,
        Amount:    amount,
        CreatedAt: time.Now(),
        Status:    "pending",
        PaymentID: "",
        Address:   address,
    }

    return purchase.ID.Hex(), s.PurchaseRepo.Create(purchase)
}

func (s *PurchaseServiceImpl) ConfirmPayment(purchaseID string) error {
    purchaseObjectID, err := bson.ObjectIDFromHex(purchaseID)
    if err != nil {
        return errors.New("invalid purchase ID format")
    }

    paymentID := generatePaymentID() // Generate a random payment ID
    return s.PurchaseRepo.UpdateStatus(purchaseObjectID, "packaging and delivering", paymentID)
}

func (s *PurchaseServiceImpl) CancelUnpaidPurchases(cutoff time.Time) (int64, error) {
    return s.PurchaseRepo.CancelOldUnpaid(cutoff)
}

func generatePaymentID() string {
    return "PAY-" + bson.NewObjectID().Hex()
}