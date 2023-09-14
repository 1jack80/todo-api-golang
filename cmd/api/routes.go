package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *api) routes() *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	// public routes
	mux.Group(func(r chi.Router) {
		r.Post("/login", a.loginHandler)
		// r.Post("signup", a.loginHandler)
		// r.Post("logout", a.logoutHandler)
	})

	// protected routes
	mux.Group(func(r chi.Router) {
		r.Use(a.sessions.Middleware)
	})
	return mux
}

func (a *api) loginHandler(w http.ResponseWriter, r *http.Request) {
	user := a.models.User.New()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.errLog.Println(err)
		return
	}

	err = a.models.User.Create(user)
	if err != nil {
		http.Error(w, "internal server err", http.StatusInternalServerError)
		a.errLog.Println(err)
		return
	}

}
