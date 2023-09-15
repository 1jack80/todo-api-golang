package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/1jack80/guardian"
	"github.com/1jack80/todo-api/models"
	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
)

func (a *api) signupHandler(w http.ResponseWriter, r *http.Request) {
	user := a.models.User.New()

	err := readJsonFromReq(r, &user)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		a.errLog.Println("unable to process json from request Body")
		return
	}
	err = a.models.User.Create(user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUsername) {
			jsonResponse(w, http.StatusBadRequest, "username already exists")
			return
		}
		jsonResponse(w, http.StatusInternalServerError, "internal server error")
		a.errLog.Println(err)
		return
	}

	jsonResponse(w, (http.StatusOK), "user created successfully")
}

func (a *api) loginHandler(w http.ResponseWriter, r *http.Request) {
	userCreds := a.models.User.New()
	readJsonFromReq(r, &userCreds)

	/*
			get usr details from request
		 	validate user details
				if invlid:
					return err
				else if valid:
					add create new session
					respond with cookie
	*/

	user, err := a.models.User.GetByUsername(userCreds.Username)
	if err != nil {
		a.errLog.Println(err)

		var mySQLError *mysql.MySQLError
		if errors.Is(err, sql.ErrNoRows) {
			jsonResponse(w, http.StatusUnauthorized, "user not found")
			return
		}
		if errors.As(err, &mySQLError) || strings.Contains(err.Error(), "sql") {
			jsonResponse(w, http.StatusInternalServerError, "internal server error")
			return
		}
		jsonResponse(w, http.StatusUnauthorized, "something went wrong")
		return
	}

	if !comparePasswordAndHash(userCreds.Password, user.Password) {
		jsonResponse(w, http.StatusUnauthorized, "incorrect password")
		return
	}

	newSession := a.sessions.CreateSession()
	newSession.Data["username"] = user.Username
	newSession.Data["userID"] = user.ID

	a.sessions.UpdateSession(newSession.ID, newSession)

	cookie, err := a.sessions.CreateCookie(newSession.ID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "Something went wrong")
	}
	http.SetCookie(w, &cookie)
	jsonResponse(w, http.StatusOK, "Login successful")
}

func (a *api) logoutHandler(w http.ResponseWriter, r *http.Request) {
	ctxVal := r.Context().Value(a.sessions.ContextKey())

	session, ok := ctxVal.(guardian.Session)
	if !ok {
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	err := a.sessions.InvalidateSession(session.ID)
	if !ok {
		a.errLog.Println(err)
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	cookie, err := a.sessions.CreateCookie(session.ID)
	if err != nil {
		a.errLog.Println(err)
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}
	cookie.Value = ""
	http.SetCookie(w, &cookie)
	jsonResponse(w, http.StatusOK, "logout successful")

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

	userID := a.getUserIDFromReqContext(r)
	if userID < 0 {
		jsonResponse(w, http.StatusBadRequest, nil)
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

func (a *api) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	strTodoID := chi.URLParam(r, "todoID")
	todoID, err := strconv.Atoi(strTodoID)
	if err != nil {
		a.errLog.Printf("todo ID: %v is not of type int", todoID)
		jsonResponse(w, http.StatusBadRequest, "Please provide a valid id of type int")
		return
	}
	userID := a.getUserIDFromReqContext(r)

	err = a.models.Todos.Delete(todoID, userID)
	if err != nil {
		a.errLog.Println(err)
		jsonResponse(w, http.StatusInternalServerError, "Unable to delete Todo")
		return
	}
	jsonResponse(w, http.StatusOK, "Todo deleted successfullyl")

}

// the request comes in with a new todo object which is used to update that of the databse
func (a *api) patchTodoHandler(w http.ResponseWriter, r *http.Request) {

	userID := a.getUserIDFromReqContext(r)

	todo := models.Todo{}
	err := readJsonFromReq(r, &todo)
	if err != nil {
		a.errLog.Printf("%v ", err)
		jsonResponse(w, http.StatusBadRequest, "Could not parse data properly")
		return
	}

	err = a.models.Todos.Patch(todo, userID)
	if err != nil {
		a.errLog.Printf("%v", err)
		jsonResponse(w, http.StatusInternalServerError, "unable to update todo")
		return
	}

	jsonResponse(w, http.StatusOK, "todo updated successfully")
}
