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

func PostCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cred models.NewCred
		readBodyErr := readBodyInStruct(r, &cred)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}
		cred.UserID = userID

		newCred, newCredErr := database.Database.NewCred(&cred)
		if newCredErr != nil {
			errorMessageResponse(w, newCredErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newCred, "application/json", http.StatusCreated)
	}
}

func GetCredList(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
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

func GetCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
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

func EditCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editCred models.NewCred
		readBodyErr := readBodyInStruct(r, &editCred)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
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

func DeleteCred(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
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
