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
	purchaseService, userService := app.InitializeServices(client) //point2<:joda bashan behtare

	// Initialize and start HTTP server
	server := http.InitializeHTTPServer(purchaseService, userService)
	server.Start()
}

/*
Here’s my program logic. Design a Camunda BPMN diagram for it. Underestand what it realy do it contains some authentication methods, a purchase process and a payment which is mocked . name each step clearly. Return the BPMN XML or a diagram I can import into Camunda Modeler.
Here is the summary presudure of the purchase process:
- Authenticated user submits a post request contains amount and address
- A put request comes into an endpoint contains purchase_id and system changes its status; if in a given time like 10 minutes the created purchase didn't paid system should canel it automatically
- END
Please return the BPMN in XML format that I can import into Camunda Modeler

*/


// p3: core: app+ domain