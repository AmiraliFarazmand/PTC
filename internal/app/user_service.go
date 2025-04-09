package app

import (
    "errors"

    "github.com/AmiraliFarazmand/PTC_Task/internal/domain"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    UserRepo domain.UserRepository
}

func (s *UserService) Signup(username, password string) error {
    // Check if the username already exists
    _, err := s.UserRepo.FindByUsername(username)
    if err == nil {
        return errors.New("username already exists")
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

func (s *UserService) Login(username, password string) (domain.User, error) {
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