package main

import (
	"log"

	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/zeebe"
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/app"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
)

func main() {
	// Initialize MongoDB
	mongoURI, _ := utils.ReadEnv("MONGO_URI")
	client := db.NewMongoDB(mongoURI)

	// Create repositories
	dbName, _ := utils.ReadEnv("DB_NAME")
	userRepo := db.NewMongoUserRepository(client.Database(dbName).Collection("Users"))
	purchaseRepo := db.NewMongoPurchaseRepository(client.Database(dbName).Collection("Purchases"))

	// Initialize services
	userService := app.InitializeUserService(userRepo)
	purchaseService := app.InitializePurchaseService(purchaseRepo)

	// Initialize Zeebe client and process manager
	zeebeClient := zeebe.NewZeebeClient()
	defer zeebe.MustCloseClient(zeebeClient)

	// Initialize process manager
	processManager := zeebe.NewZeebeProcessManager(zeebeClient)

	// Start Zeebe workers
	validateJobWorker := zeebe.ValidateCredentialsWorker(zeebeClient, userRepo)
	createUserJobWorker := zeebe.CreateUserWorker(zeebeClient, userService)
	loginCheckWorker := zeebe.CheckLoginRequestWorker(zeebeClient, userService)
	loginTokenWorker := zeebe.CreateLoginTokenWorker(zeebeClient)
	purchaseWorker := zeebe.CreatePurchaseWorker(zeebeClient, purchaseService)
	processPaymentWorker := zeebe.ProcessPaymentWorker(zeebeClient)

	defer validateJobWorker.Close()
	defer createUserJobWorker.Close()
	defer loginCheckWorker.Close()
	defer loginTokenWorker.Close()
	defer purchaseWorker.Close()
	defer processPaymentWorker.Close()

	// Initialize and start HTTP server
	server := http.NewGinServer(purchaseService, userService, processManager, zeebeClient)
	if err := server.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
