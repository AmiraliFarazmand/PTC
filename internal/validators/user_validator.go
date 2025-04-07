package validators

import (
	"errors"
	"fmt"

	"github.com/AmiraliFarazmand/PTC_Task/internal/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CheckUniquenessUsername(username string) error {
	if db.FindInstance(db.UserCollection, bson.M{"username": username}).Err() == nil {
		return fmt.Errorf("username %s already exists", username)
	}
	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 64 {
		return errors.New("username must be between 3 and 64 characters")
	}
	if err := CheckUniquenessUsername(username); err != nil {
		return err
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 64 {
		return errors.New("password must be less than 64 characters long")
	}
	return nil
}

func ValidateUsernamePassword(username, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}
