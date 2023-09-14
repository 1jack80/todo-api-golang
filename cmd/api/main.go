package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1jack80/todo-api/models"
)

type api struct {
	infoLog *log.Logger
	errLog  *log.Logger
	models  models.Models
}

func main() {
	dsn := "todo_api:todoapi_mysqldb@/TodoDB"
	addr := "127.0.0.1:8080"

	infoLog := log.New(os.Stdout, "Info ~ ", log.Ltime|log.Ldate)
	errLog := log.New(os.Stdout, "Err ~ ", log.Ltime|log.Ldate|log.Lshortfile)

	models, err := models.Init(dsn)
	if err != nil {
		errLog.Fatalf("model initialization failed: %v", err)
	}

	app := api{
		infoLog: infoLog,
		errLog:  errLog,
		models:  models,
	}

	server := http.Server{
		Addr:     addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server started on %s", addr)
	server.ListenAndServe()

}
