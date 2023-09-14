package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
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
	var mySQLError *mysql.MySQLError
	if errors.As(err, &mySQLError) {
		if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "Users.username") {
			return ErrDuplicateUsername
		}
	}
	return err
}

func (u *userModel) GetByUsername(username string) (User, error) {
	user := User{}
	stmt := "select * from Users where username=?"

	row := u.db.QueryRow(stmt, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Created, &user.Updated)
	return user, err
}
