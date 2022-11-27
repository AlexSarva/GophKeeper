package handlers

import (
	"errors"
	"net/http"
	"strings"
)

// ErrNoAuth error while valid Bearer token doesn't contain in request Header
var ErrNoAuth = errors.New("no Bearer token")

// getToken cookie selection function from Header
// returns UserID in uuid format
func getToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", ErrNoAuth
	}
	tokenValue := strings.Split(auth, "Bearer ")
	if len(tokenValue) < 2 {
		return "", ErrNoAuth
	}
	authToken := tokenValue[1]
	return authToken, nil
}
