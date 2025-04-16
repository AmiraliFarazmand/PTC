package main

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
	"github.com/AmiraliFarazmand/PTC_Task/internal/app"
)

func main() {
	// Initialize MongoDB
	client := db.NewMongoDB("mongodb://localhost:27017")
	// Create repositories
	userRepo := &db.MongoUserRepository{Collection: client.Database("ParsTasmimDB").Collection("Users")}
	purchaseRepo := &db.MongoPurchaseRepository{Collection: client.Database("ParsTasmimDB").Collection("Purchases")}

	// Initialize services
	userService := app.InitializeUserService(userRepo)
	purchaseService := app.InitializePurchaseService(purchaseRepo)

	// Initialize and start HTTP server
	server := http.InitializeHTTPServer(purchaseService, userService)
	server.Start()
}

// p3: core: app+ domain

// search for naming conventional in go
