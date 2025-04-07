package main

import (
	api "github.com/AmiraliFarazmand/PTC_Task/internal/API"
	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
)

func main() {
	db.ConnectDB("mongodb://localhost:27017", "ParsTasmimDB")
	r := api.SetupRouter()
	r.Run()

}
