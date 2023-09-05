package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (a *app) initRoutes() *chi.Mux {
	mux := chi.NewMux()
	setupMiddlewares(mux)

	mux.Post("/login", a.loginHandler)
	mux.HandleFunc("/logout", a.logoutHandler)
	mux.HandleFunc("/signup", a.signupHandler)
	mux.HandleFunc("/session_refresh", a.sessionRefreshHandler)
	mux.HandleFunc("/user_edit", a.userEditHandler)

	return mux
}

func setupMiddlewares(mux *chi.Mux) {
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	mux.Use(middleware.Logger)
}

func (a *app) loginHandler(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		a.errLog.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	// create a new session and add it to the session list
	session_cookie := a.createNewSession(creds.Username)

	http.SetCookie(w, &session_cookie)
	fmt.Fprintf(w, "login sucess, \n\n %v+ \n\n", a.sessions)
}

func (a *app) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// on logout; set expiry time of session to now
	fmt.Fprintf(w, "logout success")
}

func (a *app) signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "sign up success")
}

func (a *app) sessionRefreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Session refresh success")
}

func (a *app) userEditHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "user edited successfully")
}
