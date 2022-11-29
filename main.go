package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	sqlStmt := `
	create table if not exists todos (id integer not null primary key autoincrement, name text, description text);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	e := echo.New()

	e.Static("/", "static")
	e.File("/", "static/index.html")

	e.GET("/add", func(c echo.Context) error {
		tx, _ := db.Begin()
		stmt, _ := tx.Prepare("insert into todos(name, description) values(?,?)")
		defer stmt.Close()

		stmt.Exec("test", "description")

		tx.Commit()

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
