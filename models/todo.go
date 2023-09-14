package models

import (
	"database/sql"
	"errors"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	UserID    int       `json:"userID"`
}

var ErrNoRows = errors.New("no rows returned")

type todoModel struct {
	db *sql.DB
}

func (t *todoModel) Get(todoID int, userID int) (Todo, error) {
	todo := Todo{}
	stmt := "select * from Todos where id=? and userID=?"
	row := t.db.QueryRow(stmt, todoID, userID)
	err := row.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &todo.Created, &todo.Updated, &todo.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, ErrNoRows
		}
	}
	return todo, err
}

func (t *todoModel) GetAll(userID int) ([]*Todo, error) {
	stmt := "select * from Todos where userID=?"

	rows, err := t.db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos = []*Todo{}
	for rows.Next() {
		todo := &Todo{}

		err = rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &todo.Created, &todo.Updated, &todo.UserID)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *todoModel) Create(todo Todo) error {
	stmt := "insert into Todos(title, content, completed, userID) values (?,?,?,?)"

	_, err := t.db.Exec(stmt, todo.Title, todo.Content, todo.Completed, todo.UserID)
	return err
}
