package repository

import (
	todoapp "golang_ninja/todo-app"

	"github.com/jmoiron/sqlx"
)

type Authorizatiopn interface {
	CreateUser(user todoapp.User) (int, error)
	GetUser(username, password string) (todoapp.User, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorizatiopn
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorizatiopn: NewAuthPostgres(db),
	}
}
