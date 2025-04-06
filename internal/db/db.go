package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo"
)
var Database  mongo.Database
func ConnectDB() {

	_, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

	Database := client.Database("ParsTasmimDB")
	usersCollection :=Database.Collection("Users")
	fmt.Println(usersCollection)
	res, _ := usersCollection.InsertOne(context.Background(), bson.M{"username": "rnduser", "password": "idk"})
	id := res.InsertedID
	fmt.Println(id)
}
