package main

import (
	"log"

	todoapp "golang_ninja/todo-app"
	"golang_ninja/todo-app/pkg/handler"
)

func main() {
	srv := new(todoapp.Server)

	handlers := new(handler.Handler)

	if err := srv.Run(
		"8000",
		handlers.InitRoutes(),
	); err != nil {
		log.Fatal("error runing http server: %s", err.Error())
	}

}
