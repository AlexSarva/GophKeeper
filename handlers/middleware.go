package handlers

import (
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var ErrGetUserID = errors.New("cant get userID from ctx")

// JWTUserID type uses to pass user ID throw context
type JWTUserID string

const (
	keyPrincipalID JWTUserID = "user.id"
)

// checkContent checking content-length and content-type in basic methods of requests
func checkContent(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			headerContentType := r.Header.Get("Content-Length")
			if len(headerContentType) != 0 {
				contentLength, contentLengthErr := strconv.Atoi(headerContentType)
				if contentLengthErr != nil {
					errorMessageResponse(w, "problem in Content-Length", "application/json", http.StatusBadRequest)
					return
				}
				if contentLength != 0 {
					errorMessageResponse(w, "content-length is not equal 0", "application/json", http.StatusBadRequest)
					return
				}
			}
		}

		log.Println(r.Header)

		splittedPath := strings.Split(r.URL.Path, "/")
		lastElems := splittedPath[len(splittedPath)-2:]
		if (r.Method == "POST" || r.Method == "PATCH") && !utils.StringInSlice("files", lastElems) {
			headerContentType := r.Header.Get("Content-Type")
			if !strings.Contains("application/json, application/x-gzip", headerContentType) {
				errorMessageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// userIdentification get user-id and permissions from authorization token
func userIdentification(database *app.Storage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			jwt, jwtErr := getToken(r)
			if jwtErr != nil {
				errorMessageResponse(w, fmt.Sprint(ErrUnauthorized, ": ", jwtErr), "application/json", http.StatusUnauthorized)
				return
			}

			userID, userIDErr := database.Authorizer.ParseToken(jwt)
			if userIDErr != nil {
				errorMessageResponse(w, fmt.Sprint(ErrUnauthorized, ": ", userIDErr), "application/json", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), keyPrincipalID, userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// getUserID returns user ID from context
func getUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(keyPrincipalID).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrGetUserID
	}
	return userID, nil
}
