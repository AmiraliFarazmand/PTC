package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoPurchaseRepository struct {
	Collection *mongo.Collection
}

func NewMongoPurchaseRepository(collection *mongo.Collection) ports.PurchaseRepository {
	return &MongoPurchaseRepository{Collection: collection}
}

func (r *MongoPurchaseRepository) Create(purchase *domain.Purchase) error {  
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

func (r *MongoPurchaseRepository) GetByID(id string) (*domain.Purchase, error) {
	var purchase domain.Purchase
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.Collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&purchase)
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *MongoPurchaseRepository) Update(purchase *domain.Purchase) error {
	objectID, err := bson.ObjectIDFromHex(purchase.ID)
	if err != nil {
		return err
	}
	_, err = r.Collection.ReplaceOne(
		context.TODO(),
		bson.M{"_id": objectID},
		purchase,
	)
	return err
}

func (r *MongoPurchaseRepository) Delete(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	return err
}

func (r *MongoPurchaseRepository) GetAll() ([]*domain.Purchase, error) {
	cursor, err := r.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var purchases []*domain.Purchase
	err = cursor.All(context.TODO(), &purchases)
	return purchases, err
}

func (r *MongoPurchaseRepository) UpdateStatus(purchaseID string, status string, paymentID string, userID string) error {
	objectID, err := bson.ObjectIDFromHex(purchaseID)
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
