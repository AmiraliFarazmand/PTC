package ports

import "github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"

type UserService interface {
	Signup(username, password string) error
	Login(username, password string) (domain.User, error)
	FindUserByID(userID string) (domain.User, error)
}
