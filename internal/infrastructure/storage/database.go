package storage

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gophermart.ru/internal/entities"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(addr string) (*Database, error) {
	db := &Database{}

	if err := db.Open(addr); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) Open(addr string) error {
	conn, err := sql.Open("pgx", addr)
	if err != nil {
		return err
	}

	db.conn = conn

	if err := db.InitTable(); err != nil {
		return err
	}

	return nil
}

func (db *Database) InitTable() error {
	_, err := db.conn.Exec(
		`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL);`)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(
		`CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		number VARCHAR(255) NOT NULL UNIQUE,
		accrual BIGINT DEFAULT 0,
		id_user SERIAL NOT NULL,
		status VARCHAR(255) NOT NULL,
		date_created TIMESTAMPTZ DEFAULT NOW(),
		uploaded_at TIMESTAMPTZ DEFAULT NOW())`)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertUser(user entities.User) error {
	_, err := db.conn.Exec(
		"INSERT INTO users (name, password) VALUES ($1, $2)",
		user.Login, user.HashPassword)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) User(login string) (*entities.User, error) {
	var user entities.User
	err := db.conn.QueryRow(
		"SELECT id, name, password FROM users WHERE name = $1", login).
		Scan(&user.ID, &user.Login, &user.HashPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (db *Database) Order(number string) (*entities.Order, error) {
	var order entities.Order
	err := db.conn.QueryRow(
		"SELECT number, status, accrual, id_user, date_created, uploaded_at FROM orders WHERE number = $1", number).
		Scan(&order.Number, &order.Status, &order.Accrual, &order.IDUser, &order.CreatedAt, &order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (db *Database) Orders(IDUser int) ([]entities.Order, error) {
	var orders []entities.Order

	rows, err := db.conn.Query(
		"SELECT number, status, accrual, id_user, date_created, uploaded_at FROM orders WHERE id_user = $1 ORDER BY uploaded_at DESC", IDUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.IDUser, &order.CreatedAt, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (db *Database) AddOrder(order entities.Order) error {
	_, err := db.conn.Exec(
		"INSERT INTO orders (number, id_user, status) VALUES ($1, $2, $3)",
		order.Number, order.IDUser, order.Status)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateOrder(order entities.Order) error {
	_, err := db.conn.Exec(
		"UPDATE orders SET status = $2, accural = $3, uploaded_at = $4 WHERE number = $1",
		order.Number, order.Status, order.Accrual, time.Now())
	if err != nil {
		return err
	}

	return nil
}
