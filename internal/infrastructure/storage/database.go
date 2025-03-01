package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
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
		fmt.Println("---", err)
		return err
	}

	db.conn = conn

	// if err := b.InitTable(); err != nil {
	// 	fmt.Println("---", err)
	// 	return err
	// }
	return nil

}
