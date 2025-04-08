package purchase

import (
	"context"
	"fmt"

	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func confirmPayment(purchaseID bson.ObjectID) error {
	result, err := db.PurchaseCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": purchaseID},
		bson.M{"$set": bson.M{"status": "paid"}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no purchase found with ID: %s", purchaseID.Hex())
	}
	return nil
}
func PayPurchase(c *gin.Context) {

	var body struct {
		PurchaseID string `json:"purchase_id" binding:"required"`
	}

	// Bind JSON body to the struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON body"})
		return
	}

	purchaseID, err := bson.ObjectIDFromHex(body.PurchaseID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid purchase_id format"})
		return
	}

	// Confirm payment for the purchase
	err = confirmPayment(purchaseID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to confirm payment"})
		return
	}

	c.JSON(200, gin.H{"message": "Payment confirmed successfully"})
}
