package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *api) routes() *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	// public routes
	mux.Group(func(r chi.Router) {
		r.Post("/signup", a.signupHandler)
		r.Post("/login", a.loginHandler)
	})

	// protected routes
	mux.Group(func(r chi.Router) {
		r.Use(a.sessions.Middleware)
		r.Post("/logout", a.logoutHandler)

		r.Route("/todo", func(todoRoute chi.Router) {
			todoRoute.Get("/{todoID}", a.getOneTodoHandler)
			todoRoute.Get("/", a.getAllTodosHandler)
			todoRoute.Post("/", a.createTodoHandler)
			todoRoute.Delete("/{todoID}", a.deleteTodoHandler)
			todoRoute.Patch("/{todoID}", a.patchTodoHandler)
		})
	})
	return mux
}
