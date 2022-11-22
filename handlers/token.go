package handlers

import (
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/crypto"
	"AlexSarva/GophKeeper/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ErrNoAuth error while valid Bearer token doesn't contain in request Header
var ErrNoAuth = errors.New("no Bearer token")

// GenerateToken function of generating token for user when he successfully registered and authenticated
// based at UserID (uuid format)
// returns Token format for respond and time of expiration
func GenerateToken(userID uuid.UUID) (string, time.Time) {
	cfg := constant.GlobalContainer.Get("server-config").(models.Config)
	secret := []byte(cfg.Secret)
	session := crypto.Encrypt(userID, secret)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	return session, expiration
}

// GetToken cookie selection function from Header
// returns UserID in uuid format
func GetToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", ErrNoAuth
	}
	tokenValue := strings.Split(auth, "Bearer ")
	if len(tokenValue) < 2 {
		return "uuid.UUID{}", ErrNoAuth
	}
	authToken := tokenValue[1]
	return authToken, nil
}
