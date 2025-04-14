package app

import (
    "github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

func InitializeServices(client *mongo.Client) (PurchaseServiceImpl, UserServiceImpl) {
    purchaseCollection := client.Database("ParsTasmimDB").Collection("Purchases")
    userCollection := client.Database("ParsTasmimDB").Collection("Users")

    // Initialize repositories
    purchaseRepo := &db.MongoPurchaseRepository{Collection: purchaseCollection}
    userRepo := &db.MongoUserRepository{Collection: userCollection}

    // Initialize services
    purchaseService := PurchaseServiceImpl{PurchaseRepo: purchaseRepo}
    userService := UserServiceImpl{UserRepo: userRepo}

    return purchaseService, userService
}