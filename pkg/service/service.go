package service

import (
	todoapp "golang_ninja/todo-app"
	"golang_ninja/todo-app/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user todoapp.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type TodoList interface {
	Create(userId int, list todoapp.TodoList) (int, error)
	GetAll(userId int) ([]todoapp.TodoList, error)
	GetById(userId, listId int) (todoapp.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todoapp.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item todoapp.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todoapp.TodoItem, error)
	GetById(userId, itemId int) (todoapp.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todoapp.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorizatiopn),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemsService(repos.TodoItem, repos.TodoList),
	}
}
