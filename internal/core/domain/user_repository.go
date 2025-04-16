package domain


type User struct {
    ID       string         
    Username string         
    Password string         
}

type UserRepository interface {
    Create(user User) error
    FindByUsername(username string) (User, error)
    IsUsernameUnique(username string) (bool, error)
    FindByID(id string) (User, error)   
}
