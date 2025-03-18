package entities

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	IDUser int `json:"id_user"`
	jwt.RegisteredClaims
}

type User struct {
	ID           int    `json:"id"`
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

type Order struct {
	Number string
	IDUser int
	Status string
}

type OrdersManagment interface {
	Order(number string) (*Order, error)
	AddOrder(order Order) error
	UpdateOrder(order Order) error
}
