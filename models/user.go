package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type userModel struct {
	db *sql.DB
}

func (u *userModel) New() User {
	return User{}
}

func (u *userModel) Create(user User) error {
	stmt := "insert into Users(username, password_hash) values (?, ?)"
	binPass := sha256.Sum256([]byte(user.Password))
	strPass := hex.EncodeToString(binPass[:])

	_, err := u.db.Exec(stmt, user.Username, strPass)
	return err
}
