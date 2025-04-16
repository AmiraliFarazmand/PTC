package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoPurchaseRepository struct {
	Collection *mongo.Collection
}

func NewMongoPurchaseRepository(collection *mongo.Collection) *MongoPurchaseRepository {
	return &MongoPurchaseRepository{Collection: collection}
}

func (r *MongoPurchaseRepository) Create(purchase domain.Purchase) error {
	purchaseID, err := bson.ObjectIDFromHex(purchase.ID)
	if err != nil {
		return err
	}
	purchaseDoc := bson.M{
		"_id":        purchaseID,
		"user_id":    purchase.UserID,
		"amount":     purchase.Amount,
		"created_at": purchase.CreatedAt,
		"status":     purchase.Status,
		"payment_id": purchase.PaymentID,
		"address":    purchase.Address,
	}
	_, err = r.Collection.InsertOne(context.TODO(), purchaseDoc)
	return err
}

func (r *MongoPurchaseRepository) FindByID(id string) (domain.Purchase, error) {
	var purchase domain.Purchase
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return purchase, err
	}
	err = r.Collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&purchase)
	return purchase, err
}

func (r *MongoPurchaseRepository) UpdateStatus(id string, status string, paymentID string, userID string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	updateResult, err := r.Collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID, "user_id": userID},
		bson.M{"$set": bson.M{"status": status, "payment_id": paymentID}},
	)
	if updateResult.MatchedCount == 0 {
		return fmt.Errorf("no purchase found with this ID for current user")
	}
	return err
}

func (r *MongoPurchaseRepository) CancelOldUnpaid(cutoff time.Time) (int64, error) {
	filter := bson.M{
		"status":     "pending",
		"created_at": bson.M{"$lt": cutoff},
	}
	update := bson.M{
		"$set": bson.M{"status": "cancelled"},
	}
	result, err := r.Collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
