package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type session struct {
	session_id []byte
	username   string
	expiryTime time.Time
}

type app struct {
	errLog   *log.Logger
	infolog  *log.Logger
	sessions map[string]session
}

func main() {
	address := flag.String("addr", "127.0.0.1:8080", "Set the address for the api to listen and serve requests")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO ~ ", log.LUTC)
	errLog := log.New(os.Stderr, "ERROR ~ ", log.LUTC|log.Lshortfile)

	api := &app{
		errLog:   errLog,
		infolog:  infoLog,
		sessions: map[string]session{},
	}

	server := &http.Server{
		Addr:     *address,
		ErrorLog: errLog,
		Handler:  api.initRoutes(),
	}

	infoLog.Printf("Server started on %s", server.Addr)
	errLog.Fatal(server.ListenAndServe())
}
