package db

import (
    "context"
    "time"

    "github.com/AmiraliFarazmand/PTC_Task/internal/domain"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoPurchaseRepository struct {
    Collection *mongo.Collection
}

func (r *MongoPurchaseRepository) Create(purchase domain.Purchase) error {
    _, err := r.Collection.InsertOne(context.TODO(), purchase)
    return err
}

func (r *MongoPurchaseRepository) FindByID(id bson.ObjectID) (domain.Purchase, error) {
    var purchase domain.Purchase
    err := r.Collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&purchase)
    return purchase, err
}

func (r *MongoPurchaseRepository) UpdateStatus(id bson.ObjectID, status string, paymentID string) error {
    _, err := r.Collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": id},
        bson.M{"$set": bson.M{"status": status, "payment_id": paymentID}},
    )
    return err
}

func (r *MongoPurchaseRepository) CancelOldUnpaid(cutoff time.Time) (int64, error) {
    filter := bson.M{
        "status":    "pending",
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