package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func ConnectDB(uri, dbName string) {
	_, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, _ := mongo.Connect(options.Client().ApplyURI(uri))
	DB = client.Database(dbName)
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

func InsertIntoCollection(collection *mongo.Collection, instance bson.M) {
	res, _ := collection.InsertOne(context.Background(), instance)
	id := res.InsertedID
	fmt.Println(id)
}

func FindInstance(collection *mongo.Collection, instance bson.M) *mongo.SingleResult {
	return collection.FindOne(context.Background(), instance)
}

func DeleteInstance(collection *mongo.Collection, instance bson.M) {
	collection.DeleteOne(context.Background(), instance)
}

func UpdateInstance(collection *mongo.Collection, filter bson.M, update bson.M) {
	collection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
}

