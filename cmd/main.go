package main

import (
	api "github.com/AmiraliFarazmand/PTC_Task/internal/API"
	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"github.com/AmiraliFarazmand/PTC_Task/internal/purchase"
)

func main() {

	go purchase.CancelUnpaidPurchases()
	db.ConnectDB("mongodb://localhost:27017", "ParsTasmimDB")
	r := api.SetupRouter()
	r.Run()

}
