package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
    ID       bson.ObjectID `bson:"_id,omitempty"`
    Username string             `bson:"username"`
    Password string             `bson:"password"`
}

type UserRepository interface {
    Create(user User) error
    FindByUsername(username string) (User, error)
    // IsUsernameUnique(username string) (bool, error)
}
