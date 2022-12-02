package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"todo/controller"
	"todo/service"
	"todo/storage"
)

func main() {
	todoStorage := initStorage()
	defer func(todoStorage storage.TodoStorage) {
		err := todoStorage.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(todoStorage)

	s := service.NewTodoService(todoStorage)

	e := echo.New()
	controller.NewTodoController(e, s).RegisterRoutes()

	e.Logger.Fatal(e.Start(":8080"))
}

func initStorage() storage.TodoStorage {
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}

	s := storage.NewTodoStorage(db)
	err = s.Init()
	if err != nil {
		panic(err)
	}

	return s
}
