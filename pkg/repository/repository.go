package repository

import "github.com/jmoiron/sqlx"

type Authorizatiopn interface {
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
	return &Repository{}
}
