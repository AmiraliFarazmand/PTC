package main

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/adapters/http"
	"github.com/AmiraliFarazmand/PTC_Task/internal/app"
)

func main() {
	// Initialize MongoDB
	client := db.InitializeMongoDB("mongodb://localhost:27017")

	// Initialize services
	purchaseService, userService := app.InitializeServices(client)

	// Initialize and start HTTP server
	server := http.InitializeHTTPServer(purchaseService, userService)
	server.Start()
}
