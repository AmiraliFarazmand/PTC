package domain

import (
    "time"

    "go.mongodb.org/mongo-driver/v2/bson"
)

type Purchase struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    UserID    bson.ObjectID `bson:"user_id"`
    Amount    int                `bson:"amount"`
    CreatedAt time.Time          `bson:"created_at"`
    Status    string             `bson:"status"`
    PaymentID string             `bson:"payment_id"`
    Address   string             `bson:"address"`
}

type PurchaseRepository interface {
    Create(purchase Purchase) error
    FindByID(id bson.ObjectID) (Purchase, error)
    UpdateStatus(id bson.ObjectID, status string, paymentID string) error
    CancelOldUnpaid(cutoff time.Time) (int64, error)
}