package main

import (
    "time"

    "github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
    "github.com/AmiraliFarazmand/PTC_Task/internal/app"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
    // Connect to MongoDB
    client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        panic(err)
    }
    purchaseCollection := client.Database("ParsTasmimDB").Collection("Purchases")

    // Initialize repositories
    purchaseRepo := &db.MongoPurchaseRepository{Collection: purchaseCollection}

    // Initialize services
    purchaseService := &app.PurchaseServiceImpl{PurchaseRepo: purchaseRepo}

    // Setup routes
    r := gin.Default()
    r.POST("/purchase", func(c *gin.Context) {
        var body struct {
            UserID  string `json:"user_id" binding:"required"`
            Amount  int    `json:"amount" binding:"required"`
            Address string `json:"address" binding:"required"`
        }
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request body"})
            return
        }
        err := purchaseService.CreatePurchase(body.UserID, body.Amount, body.Address)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(201, gin.H{"message": "Purchase created successfully"})
    })

    r.POST("/purchase/pay", func(c *gin.Context) {
        var body struct {
            PurchaseID string `json:"purchase_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request body"})
            return
        }
        err := purchaseService.ConfirmPayment(body.PurchaseID)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"message": "Payment confirmed successfully"})
    })

    go func() {
        for {
            time.Sleep(15 * time.Second)
            _, err := purchaseService.CancelUnpaidPurchases(time.Now().Add(-1 * time.Minute))
            if err != nil {
                // Log the error
            }
        }
    }()

    // Start the server
    r.Run(":8080")
}