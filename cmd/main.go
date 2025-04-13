package main

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
	"github.com/AmiraliFarazmand/PTC_Task/internal/app"
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
	userCollection := client.Database("ParsTasmimDB").Collection("Users")

	// Initialize repositories
	purchaseRepo := &db.MongoPurchaseRepository{Collection: purchaseCollection}
	userRepo := &db.MongoUserRepository{Collection: userCollection}

	// Initialize services
	purchaseService := app.PurchaseServiceImpl{PurchaseRepo: purchaseRepo}
	userService := app.UserServiceImpl{UserRepo: userRepo}

	// Initialize HTTP server
	server := http.NewGinServer(purchaseService, userService)

	// Start the server
	server.Start()
}
