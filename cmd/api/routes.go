package main

import (
	"github.com/go-chi/chi"
)

func (a *api) routes() *chi.Mux {
	mux := chi.NewMux()

	return mux
}
