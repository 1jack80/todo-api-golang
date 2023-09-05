package main

import (
	"crypto/sha256"
	"time"
)

func (a *app) createNewSession(username string) {

	session_id := string(sha256.New().Sum([]byte(username)))

	a.sessions[session_id] = session{
		username:   username,
		expiryTime: time.Now().Add(3 * time.Minute),
	}
}
