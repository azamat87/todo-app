package main

import (
	todoapp "golang_ninja/todo-app"
	"golang_ninja/todo-app/pkg/handler"
	"golang_ninja/todo-app/pkg/repository"
	"golang_ninja/todo-app/pkg/service"
	"log"

	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error error instal configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("failed init db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todoapp.Server)

	if err := srv.Run(
		viper.GetString("port"),
		handlers.InitRoutes(),
	); err != nil {
		log.Fatalf("error runing http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
