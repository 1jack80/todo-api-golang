package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/1jack80/todo-api/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *api) routes() *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	// public routes
	mux.Group(func(r chi.Router) {
		r.Post("/signup", a.signupHandler)
		// r.Post("signup", a.loginHandler)
		// r.Post("logout", a.logoutHandler)
	})

	// protected routes
	mux.Group(func(r chi.Router) {
		r.Use(a.sessions.Middleware)
	})
	return mux
}

func (a *api) signupHandler(w http.ResponseWriter, r *http.Request) {
	user := a.models.User.New()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.errLog.Println(err)
		return
	}

	err = a.models.User.Create(user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUsername) {
			http.Error(w, "username not available", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server err", http.StatusInternalServerError)
		a.errLog.Println(err)
		return
	}

	// todo: send a json response: status: ok , msg user created successfully

}
