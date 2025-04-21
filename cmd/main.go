package main

import (

	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/zeebe"
	"github.com/AmiraliFarazmand/PTC_Task/internal/utils"
	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/worker"
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/app"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
)

func main() {
	// Initialize MongoDB
	mongoURI, _ := utils.ReadEnv("MONGO_URI")
	client := db.NewMongoDB(mongoURI)

	// Create repositories
	dbName, _ := utils.ReadEnv("DB_NAME")
	userRepo := db.NewMongoUserRepository(client.Database(dbName).Collection("Users"))
	purchaseRepo := db.NewMongoPurchaseRepository(client.Database("ParsTasmimDB").Collection("Purchases"))

	// Initialize services
	userService:= app.InitializeUserService(userRepo)
	purchaseService := app.InitializePurchaseService(purchaseRepo)
	zeebeClient := zeebe.NewZeebeClient()
	defer zeebe.MustCloseClient(zeebeClient)
	// Deploy BPMN process
	zeebe.DeploySignupProcess(zeebeClient)
	// Start workers
	var validateJobWorker, createUserJobWorker worker.JobWorker

	validateJobWorker = zeebe.ValidateCredentialsWorker(zeebeClient)
	createUserJobWorker = zeebe.CreateUserWorker(zeebeClient, userRepo)
	defer validateJobWorker.Close()
	defer createUserJobWorker.Close()
	// Initialize and start HTTP server
	server := http.InitializeHTTPServer(purchaseService, userService)
	server.Start()
	select {}
}

// search for naming conventional in go
