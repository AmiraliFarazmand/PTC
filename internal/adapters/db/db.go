package db
//TODO: rename
import (
	"context"
	"log"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{Collection: collection}
}
func (r *MongoUserRepository) Create(user domain.User) error {
	objectID := bson.NewObjectID()
	_, err := r.Collection.InsertOne(context.TODO(), bson.M{
		"_id":      objectID,
		"username": user.Username,
		"password": user.Password,
	})
	return err
}
func (r *MongoUserRepository) FindByUsername(username string) (domain.User, error) {
	var result struct {
		ID       bson.ObjectID `bson:"_id"`
		Username string        `bson:"username"`
		Password string        `bson:"password"`
	}
	var user domain.User
	err := r.Collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		return user, err
	}
	user.ID = result.ID.Hex()
	user.Username = result.Username
	user.Password = result.Password
	return user, err
}

func (r *MongoUserRepository) FindByID(id string) (domain.User, error) {
	var result struct {
		ID       bson.ObjectID `bson:"_id"`
		Username string        `bson:"username"`
		Password string        `bson:"password"`
	}
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	err = r.Collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
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

func NewMongoDB(uri string) *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client
}
