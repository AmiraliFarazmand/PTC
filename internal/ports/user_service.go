package ports

import "github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"

type UserService interface {
	Signup(username, password string) error
	Login(username, password string) (domain.User, error)
	FindUserByID(userID string) (domain.User, error)
}

type ZeebeProcessManager interface {
	StartSignupProcess(username, password string) error
	StartLoginProcess(username, password string) error
}
