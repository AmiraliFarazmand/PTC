package app

import (
    "github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
    "go.mongodb.org/mongo-driver/v2/mongo"  //adaptor bayad bashe na mongo
)

func InitializePurchaseService(client *mongo.Client) PurchaseServiceImpl {
    purchaseCollection := client.Database("ParsTasmimDB").Collection("Purchases")
    purchaseRepo := &db.MongoPurchaseRepository{Collection: purchaseCollection}
    return PurchaseServiceImpl{PurchaseRepo: purchaseRepo}
}

func InitializeUserService(client *mongo.Client) UserServiceImpl {
    userCollection := client.Database("ParsTasmimDB").Collection("Users")
    userRepo := &db.MongoUserRepository{Collection: userCollection}
    return UserServiceImpl{UserRepo: userRepo}
}