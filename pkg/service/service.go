package service

import "golang_ninja/todo-app/pkg/repository"

type Authorizatiopn interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorizatiopn
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
