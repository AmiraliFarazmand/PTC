package main

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/app"
)

func main() {
	// Initialize MongoDB
	client := db.NewMongoDB("mongodb://localhost:27017")

	// Create repositories
	userRepo := db.NewMongoUserRepository(client.Database("ParsTasmimDB").Collection("Users"))
	purchaseRepo := db.NewMongoPurchaseRepository(client.Database("ParsTasmimDB").Collection("Purchases"))

	// Initialize services
	userService := app.InitializeUserService(userRepo)
	purchaseService := app.InitializePurchaseService(purchaseRepo)

	// Initialize and start HTTP server
	server := http.InitializeHTTPServer(purchaseService, userService)
	server.Start()
}


// search for naming conventional in go
