package storage

import (
	"database/sql"

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
		`CREATE TABLE IF NOT EXISTS status_orders (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255)  NOT NULL UNIQUE);`)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(
		`CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		number VARCHAR(255) NOT NULL UNIQUE,
		id_user SERIAL NOT NULL,
		status VARCHAR(255) NOT NULL);`)
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
		Scan(&user.Id, &user.Login, &user.HashPassword)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) Order(number string) (*entities.Order, error) {
	var order entities.Order
	err := db.conn.QueryRow(
		"SELECT number, id_user, status FROM orders WHERE number = $1", number).
		Scan(&order.Number, &order.IdUser, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (db *Database) AddOrder(order entities.Order) error {
	_, err := db.conn.Exec(
		"INSERT INTO orders (number, id_user, status) VALUES ($1, $2, $3)",
		order.Number, order.IdUser, order.Status)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateOrder(order entities.Order) error {
	_, err := db.conn.Exec(
		"UPDATE orders SET status = $2 WHERE number = $1",
		order.Number, order.Status)
	if err != nil {
		return err
	}

	return nil
}
