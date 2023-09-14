package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
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
