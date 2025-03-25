package entities

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoMoney        = errors.New("на счету недостаточно средств")
	ErrIncorrectOrder = errors.New("неверный номер заказа")
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
	OrderNew        string = "NEW"
	OrderProcessing string = "PROCESSING"
	OrderInvalid    string = "INVALID"
	OrderProcessed  string = "PROCESSED"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
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

type Wallet struct {
	ID       int     `json:"-"`
	Balance  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}

type Withdraw struct {
	ID        int       `json:"-"`
	Order     string    `json:"order"`
	Sum       float64   `json:"sum"`
	CreatedAt time.Time `json:"processed_at,omitempty"`
}

type WalletManagment interface {
	Balance(IDUser int) (*Wallet, error)
	Withdraw(w Withdraw) error
	SumWithdraw(IDUser int) (float64, error)
	Withdraws(IDUser int) ([]Withdraw, error)
}
