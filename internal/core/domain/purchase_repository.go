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

type PurchaseRepository interface {
    Create(purchase Purchase) error
    FindByID(id string) (Purchase, error)
    UpdateStatus(id string, status string, paymentID string, userID string) error
    CancelOldUnpaid(cutoff time.Time) (int64, error)
}