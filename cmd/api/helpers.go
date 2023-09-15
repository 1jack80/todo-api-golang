package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/1jack80/guardian"
)

// decodes the json object into the object
// using this function requres that you pass a pointer as the object parameter
func readJsonFromReq(r *http.Request, objPointer interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&objPointer)
	return err
}

func jsonResponse(w http.ResponseWriter, statusCode int, msg interface{}) {
	response := struct {
		Status string      `json:"status"`
		Msg    interface{} `json:"msg"`
	}{Status: http.StatusText(statusCode), Msg: msg}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func comparePasswordAndHash(password string, hash string) bool {
	binPass := sha256.Sum256([]byte(password))
	strHash := hex.EncodeToString(binPass[:])

	return (strHash == hash)
}

func (a *api) getUserIDFromReqContext(r *http.Request) int {
	ctxVal := r.Context().Value(a.sessions.ContextKey())
	session, ok := ctxVal.(guardian.Session)
	if !ok {
		a.errLog.Printf("unable to get session data\n")
		return -1
	}
	userIDstr, ok := session.Data["userID"]
	if !ok {
		a.errLog.Printf("user id was not found in session data\n")
		return -1
	}

	userID, ok := userIDstr.(int)
	if !ok {
		a.errLog.Printf("usersID is not of type int")
		return -1
	}
	return userID
}
