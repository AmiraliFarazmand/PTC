package app

import (
	"errors"

	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
	"github.com/AmiraliFarazmand/PTC_Task/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepo ports.UserRepository
}

func (s *UserServiceImpl) Signup(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return s.UserRepo.Create(user)
}

func (s *UserServiceImpl) Login(username, password string) (ports.UserDTO, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return ports.UserDTO{}, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return ports.UserDTO{}, errors.New("invalid username or password")
	}

	return ports.UserDTO{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *UserServiceImpl) FindUserByID(userID string) (ports.UserDTO, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return ports.UserDTO{}, errors.New("user not found")
	}

	return ports.UserDTO{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
