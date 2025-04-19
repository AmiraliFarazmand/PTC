package main

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/app"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
)

func main() {
	// Initialize MongoDB
	mongoURI,_ := utils.ReadEnv("MONGO_URI")
	client := db.NewMongoDB(mongoURI)

	// Create repositories
	dbName, _ := utils.ReadEnv("DB_NAME")
	userRepo := db.NewMongoUserRepository(client.Database(dbName).Collection("Users"))
	purchaseRepo := db.NewMongoPurchaseRepository(client.Database("ParsTasmimDB").Collection("Purchases"))

	// Initialize services
	userService := app.InitializeUserService(userRepo)
	purchaseService := app.InitializePurchaseService(purchaseRepo)

	// Initialize and start HTTP server
	server := http.InitializeHTTPServer(purchaseService, userService)
	server.Start()
}


// search for naming conventional in go
