package handlers

import (
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
			jwt, jwtErr := GetToken(r)
			if jwtErr != nil {
				errorMessageResponse(w, fmt.Sprint(ErrUnauthorized, ": ", jwtErr), "application/json", http.StatusUnauthorized)
				return
			}

			userID, userIDErr := database.Authorizer.ParseToken(jwt)
			if userIDErr != nil {
				errorMessageResponse(w, fmt.Sprint(ErrUnauthorized, ": ", userIDErr), "application/json", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user.id", userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
