package entities

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	HashPassword string `json:"password"`
}

func (u *User) IsEcual(user *User) bool {
	return u.Login == user.Login && u.HashPassword == user.HashPassword
}

type UserManagment interface {
	InsertUser(User) error
	User(login string) (*User, error)
}
