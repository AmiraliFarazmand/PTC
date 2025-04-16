package app

import (
	"errors"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepo domain.UserRepository
}

func (s *UserServiceImpl) Signup(username, password string) error {
	// Check if the username already exists
	_, err := s.UserRepo.FindByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	if err = validateUsername(username); err != nil {
		return err
	}
	if err = validatePassword(password); err != nil {
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create the user
	user := domain.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return s.UserRepo.Create(user)
}

func (s *UserServiceImpl) Login(username, password string) (domain.User, error) {
	// Find the user by username
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return domain.User{}, errors.New("invalid username or password")
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, errors.New("invalid username or password")
	}

	return user, nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 64 {
		return errors.New("username must be between 3 and 64 characters")
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

func (s *UserServiceImpl) FindUserByID(userID string) (domain.User, error) {
	// Use the repository to find the user
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}
