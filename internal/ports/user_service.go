package ports

type UserService interface {
	Signup(username, password string) error
	Login(username, password string) (UserDTO, error)
	FindUserByID(userID string) (UserDTO, error)
}
