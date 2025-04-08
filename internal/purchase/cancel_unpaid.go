package purchase

import (
	"context"
	"log"
	"time"

	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CancelUnpaidPurchases() {
	ticker := time.NewTicker(15 * time.Second) 
    defer ticker.Stop()
    
    for range ticker.C {
        cutoff := time.Now().Add(-1 * time.Minute)

        filter := bson.M{
            "status":    "pending",
            "created_at": bson.M{"$lt": cutoff},
        }

        update := bson.M{
            "$set": bson.M{"status": "cancelled"},
        }

        result, err := db.PurchaseCollection.UpdateMany(context.TODO(), filter, update)
        if err != nil {
            log.Println("Auto-cancel failed:", err)
        } else if result.ModifiedCount > 0 {
            log.Printf("Auto-cancelled %d old purchases\n", result.ModifiedCount)
        }
    }
}           