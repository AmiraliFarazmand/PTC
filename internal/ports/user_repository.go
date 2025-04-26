package ports

import (
	"github.com/AmiraliFarazmand/PTC_Task/internal/core/domain"
)

type UserRepository interface {
	Create(user domain.User) error
	FindByUsername(username string) (domain.User, error)
	IsUsernameUnique(username string) (bool, error)
	FindByID(id string) (domain.User, error)
}
