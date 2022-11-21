package handlers

import (
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/storage"
	"context"
	"errors"
	"fmt"
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
		if r.Method == "POST" || r.Method == "PATCH" {
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
func userIdentification(database *app.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userID, tokenErr := GetToken(r)
			if tokenErr != nil {
				errorMessageResponse(w, fmt.Sprint(ErrUnauthorized, ": ", tokenErr), "application/json", http.StatusUnauthorized)
				return
			}

			if !database.Admin.CheckUser(userID) {
				errorMessageResponse(w, "user doesnt registered. visit: api/v1/register", "application/json", http.StatusUnauthorized)
				return
			}

			userRoles, userRolesErr := database.Admin.GetUserRoles(userID)
			if userRolesErr != nil {
				if errors.Is(storage.ErrNoValues, userRolesErr) {
					errorMessageResponse(w, ErrUnauthorized.Error()+": user doesnt have any role", "application/json", http.StatusUnauthorized)
					return
				}
				errorMessageResponse(w, userRolesErr.Error(), "application/json", http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), "acl.permission", userRoles)
			ctx = context.WithValue(ctx, "userID", userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
