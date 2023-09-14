package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Models struct {
	db *sql.DB
}

// connect to the database using the provided dsn
func Init(dsn string) (Models, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return Models{}, err
	}
	return Models{db: db}, err
}
