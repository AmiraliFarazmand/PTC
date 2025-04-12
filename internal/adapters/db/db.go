package db

import (
    "context"

    "github.com/AmiraliFarazmand/PTC_Task/internal/domain"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoUserRepository struct {
    Collection *mongo.Collection
}

func (r *MongoUserRepository) Create(user domain.User) error {
    _, err := r.Collection.InsertOne(context.TODO(), user)
    return err
}

func (r *MongoUserRepository) FindByUsername(username string) (domain.User, error) {
    var user domain.User
    err := r.Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
    return user, err
}
func (r *MongoUserRepository) FindByID(id string) (domain.User, error) {
    var user domain.User
    err := r.Collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
    return user, err
}

func (r *MongoUserRepository) IsUsernameUnique(username string) (bool, error) {
    filter := bson.M{"username": username}
    count, err := r.Collection.CountDocuments(context.TODO(), filter)
    if err != nil {
        return false, err
    }
    return count == 0, nil
}