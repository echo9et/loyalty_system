package entities

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

const (
	ORDER_NEW        string = "NEW"
	ORDER_PROCESSING string = "PROCESSING"
	ORDER_INVALID    string = "INVALID"
	ORDER_PROCESSED  string = "PROCESSED"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	IDUser     int       `json:"-"`
	CreatedAt  time.Time `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type OrdersManagment interface {
	Order(number string) (*Order, error)
	Orders(IDUser int) ([]Order, error)
	AddOrder(order Order) error
	UpdateOrder(order Order) error
}
