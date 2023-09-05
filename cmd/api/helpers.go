package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func (a *app) createNewSession(username string) http.Cookie {
	hour := time.Now().UTC().Minute()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d%s", hour, username)))

	session_id := hex.EncodeToString(hash[:])

	expiryTime := time.Now().Add(3 * time.Minute)

	a.sessions[session_id] = session{
		username:   username,
		expiryTime: expiryTime,
	}

	return http.Cookie{
		Name:    "session_token",
		Value:   session_id,
		Expires: expiryTime,
	}
}
