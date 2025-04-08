package main

import (
	"sync"

	api "github.com/AmiraliFarazmand/PTC_Task/internal/API"
	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/purchase"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		purchase.CancelUnpaidPurchases()
	}()

	db.ConnectDB("mongodb://localhost:27017", "ParsTasmimDB")
	r := api.SetupRouter()
	r.Run()

	wg.Wait() // Wait for the goroutine to finish
}
