package db

import (
	"context"
	"log"

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
	log.Println(user, err)
	return user, err
}

func (r *MongoUserRepository) FindByID(id string) (domain.User, error) {
	var result struct {
		ID       bson.ObjectID `bson:"_id"`
		Username string        `bson:"username"`
		Password string 		`bson:"password"`
	}
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID format:", err, id)
		return domain.User{}, err
	}

	err = r.Collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		log.Println("In FindByID method", err)
		return domain.User{}, err
	}

	user := domain.User{
		ID:       result.ID.Hex(),
		Username: result.Username,
		Password: result.Password,
	}
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
