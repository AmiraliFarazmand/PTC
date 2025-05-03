package main

import (
	"log"

	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/zeebe"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/zeebe/workers"
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
	validateJobWorker := workers.ValidateCredentialsWorker(zeebeClient, userRepo)
	createUserJobWorker := workers.CreateUserWorker(zeebeClient, userService)
	loginCheckWorker := workers.CheckLoginRequestWorker(zeebeClient, userService)
	loginTokenWorker := workers.CreateLoginTokenWorker(zeebeClient)
	purchaseWorker := workers.CreatePurchaseWorker(zeebeClient, purchaseService)
	cancelUnpaidPurchaeWorker := workers.CancelUnpaidPurchaseWorker(zeebeClient, purchaseService)
	processPaymentWorker := workers.ProcessPaymentWorker(zeebeClient)

	defer validateJobWorker.Close()
	defer createUserJobWorker.Close()
	defer loginCheckWorker.Close()
	defer loginTokenWorker.Close()
	defer purchaseWorker.Close()
	defer cancelUnpaidPurchaeWorker.Close()
	defer processPaymentWorker.Close()

	// Initialize and start HTTP server
	server := http.NewGinServer(purchaseService, userService, processManager, zeebeClient)
	if err := server.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
