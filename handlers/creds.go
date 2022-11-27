package handlers

import (
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// PostCred - add credential method
//
// Handler POST /api/v1/info/creds
//
//	"title": "<title>",
//	"login": "<login>",
//	"password": "<password>",
//	"notes": "<notes>"
//
// Possible response codes:
// 201 - credential successfully added;
// 400 - invalid request format;
// 401 - problem from authentication;
// 500 - an internal server error.
func PostCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cred models.NewCred
		readBodyErr := readBodyInStruct(r, &cred)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		userID, userIDErr := getUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}
		cred.UserID = userID

		if cred.Login == "" || cred.Passwd == "" {
			errorMessageResponse(w, "empty fields error", "application/json", http.StatusBadRequest)
			return
		}

		newCred, newCredErr := database.Database.NewCred(&cred)
		if newCredErr != nil {
			errorMessageResponse(w, newCredErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newCred, "application/json", http.StatusCreated)
	}
}

// GetCredList - get all credentials method
//
// Handler GET /api/v1/info/creds
//
// Possible response codes:
// 200 - returns information;
// 204 - no values in database;
// 400 - invalid request format;
// 401 - problem from authentication;
// 500 - an internal server error.
func GetCredList(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := getUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		creds, credsErr := database.Database.AllCreds(userID)
		if credsErr != nil {
			errorMessageResponse(w, credsErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		if len(creds) == 0 {
			errorMessageResponse(w, "no values", "application/json", http.StatusNoContent)
			return
		}

		resultResponse(w, creds, "application/json", http.StatusOK)
	}
}

// GetCred - get credential method (by uuid)
//
// Handler GET /api/v1/info/creds/{id}
//
// Possible response codes:
// 200 - returns information;
// 204 - no values in database;
// 400 - invalid request format;
// 401 - problem from authentication;
// 409 - no such credential in database;
// 500 - an internal server error.
func GetCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := getUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		credIDStr := chi.URLParam(r, "id")
		credUUID, credUUIDErr := uuid.Parse(credIDStr)
		if credUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		cred, credErr := database.Database.GetCred(credUUID, userID)
		if credErr != nil {
			if errors.Is(credErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, credErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		resultResponse(w, cred, "application/json", http.StatusOK)
	}
}

// EditCred - edit credential information method
//
// Handler PATCH /api/v1/info/creds/{id}
//
//	"title": "<title>",
//	"login": "<login>",
//	"password": "<password>",
//	"notes": "<notes>"
//
// Possible response codes:
// 201 - credential information successfully changed;
// 400 - invalid request format;
// 401 - problem from authentication;
// 409 - no such credential in database;
// 500 - an internal server error.
func EditCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editCred models.NewCred
		readBodyErr := readBodyInStruct(r, &editCred)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		userID, userIDErr := getUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		credIDStr := chi.URLParam(r, "id")
		credUUID, credUUIDErr := uuid.Parse(credIDStr)
		if credUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		cred, credErr := database.Database.GetCred(credUUID, userID)
		if credErr != nil {
			if errors.Is(credErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such cred in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, credErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if editCred.Title == "" {
			editCred.Title = cred.Title
		}

		if editCred.Login == "" {
			editCred.Login = cred.Login
		}

		if editCred.Passwd == "" {
			editCred.Passwd = cred.Passwd
		}

		if editCred.Notes == "" {
			editCred.Notes = cred.Notes
		}

		editCred.ID = cred.ID
		editCred.UserID = userID

		newCred, newCredErr := database.Database.EditCred(editCred)
		if newCredErr != nil {
			if errors.Is(newCredErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such cred in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, newCredErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newCred, "application/json", http.StatusCreated)
	}
}

// DeleteCred - delete credential method
//
// Handler DELETE /api/v1/info/creds/{id}
//
// Possible response codes:
// 200 - successful deleted;
// 400 - invalid request format;
// 401 - problem from authentication;
// 409 - no such credential in database;
// 500 - an internal server error.
func DeleteCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := getUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		credIDStr := chi.URLParam(r, "id")
		credUUID, credUUIDErr := uuid.Parse(credIDStr)
		if credUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		delErr := database.Database.DeleteCred(credUUID, userID)
		if delErr != nil {
			if errors.Is(delErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}
			errorMessageResponse(w, delErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, "successful deleted", "application/json", http.StatusOK)
	}
}
