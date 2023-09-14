package main

import (
	"net/http"
	"strconv"

	"github.com/1jack80/guardian"
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

		})
	})
	return mux
}

func (a *api) getOneTodoHandler(w http.ResponseWriter, r *http.Request) {
	/*
		get user from context
		use userid to get all the requested todo by user
	*/
	ctx := r.Context().Value(a.sessions.ContextKey())

	session, ok := ctx.(guardian.Session)
	if !ok {
		jsonResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	userIDStr, ok := session.Data["userID"]
	if !ok {
		a.errLog.Println("unable to get userId from sessin data")
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}
	userID, ok := userIDStr.(int)
	if !ok {
		a.errLog.Println("user id obtained from session data is not of type int")
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	todoIDstr := chi.URLParam(r, "todoID")
	todoID, err := strconv.Atoi(todoIDstr)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "invalid todo ID")
		return
	}

	todo, err := a.models.Todos.Get(todoID, userID)
	if err != nil {
		a.errLog.Println(err)
		jsonResponse(w, http.StatusBadRequest, "could not get todo")
		return
	}

	jsonResponse(w, http.StatusOK, todo)
}

func (a *api) getAllTodosHandler(w http.ResponseWriter, r *http.Request) {
	ctxVal := r.Context().Value(a.sessions.ContextKey())
	session, ok := ctxVal.(guardian.Session)
	if !ok {
		a.errLog.Printf("unable to get session data\n")
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}
	userIDstr, ok := session.Data["userID"]
	if !ok {
		a.errLog.Printf("user id was not found in session data\n")
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	userID, ok := userIDstr.(int)
	if !ok {
		a.errLog.Printf("usersID is not of type int")
		jsonResponse(w, http.StatusBadRequest, nil)
		return
	}

	todos, err := a.models.Todos.GetAll(userID)
	if err != nil {
		a.errLog.Printf("unable to get all todos: %v ", err)
		jsonResponse(w, http.StatusBadRequest, nil)
		return
	}

	jsonResponse(w, http.StatusOK, todos)
}

func (a *api) createTodoHandler(w http.ResponseWriter, r *http.Request) {

	ctxVal := r.Context().Value(a.sessions.ContextKey())
	session, ok := ctxVal.(guardian.Session)
	if !ok {
		a.errLog.Printf("unable to get session data\n")
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}
	userIDstr, ok := session.Data["userID"]
	if !ok {
		a.errLog.Printf("user id was not found in session data\n")
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	userID, ok := userIDstr.(int)
	if !ok {
		a.errLog.Printf("usersID is not of type int")
		jsonResponse(w, http.StatusBadRequest, nil)
		return
	}

	todo := models.Todo{}
	err := readJsonFromReq(r, &todo)
	if err != nil {
		a.errLog.Printf("unable to read todos from request body:\n %v", err)
		jsonResponse(w, http.StatusBadRequest, nil)
		return
	}

	todo.UserID = userID
	a.infoLog.Printf("%+v", todo)

	err = a.models.Todos.Create(todo)
	if err != nil {
		a.errLog.Printf("unable to create Todo %v", err)
		jsonResponse(w, http.StatusInternalServerError, nil)
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}
