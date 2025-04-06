package main

import (
	"fmt"

	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	db.ConnectDB("mongodb://localhost:27017", "ParsTasmimDB")
	userCollection := db.GetCollection("Users")
	// db.InsertIntoCollection(userCollection, bson.M{"username": "sdfsf", "password": "sdfsdf"})
	found := db.FindInstance(userCollection, bson.M{"username":"sdfsf"})
	fmt.Println(found)
	
}
