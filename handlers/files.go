package handlers

import (
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"AlexSarva/GophKeeper/utils"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func PostFile(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var file models.NewFile
		readBodyErr := readBodyInStruct(r, &file)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}
		file.UserID = userID

		newFile, newFileErr := database.Database.NewFile(&file)
		if newFileErr != nil {
			errorMessageResponse(w, newFileErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newFile, "application/json", http.StatusCreated)
	}
}

func GetFileList(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		files, filesErr := database.Database.AllFiles(userID)
		if filesErr != nil {
			errorMessageResponse(w, filesErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		if len(files) == 0 {
			errorMessageResponse(w, "no values", "application/json", http.StatusNoContent)
			return
		}

		resultResponse(w, files, "application/json", http.StatusOK)
	}
}

func GetFile(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		fileIDStr := chi.URLParam(r, "id")
		fileUUID, fileUUIDErr := uuid.Parse(fileIDStr)
		if fileUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		file, fileErr := database.Database.GetFile(fileUUID, userID)
		if fileErr != nil {
			if errors.As(fileErr, &storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, fileErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		resultResponse(w, file, "application/json", http.StatusOK)
	}
}

func EditFile(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editFile models.NewFile
		readBodyErr := readBodyInStruct(r, &editFile)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		fileIDStr := chi.URLParam(r, "id")
		fileUUID, fileUUIDErr := uuid.Parse(fileIDStr)
		if fileUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		file, fileErr := database.Database.GetFile(fileUUID, userID)
		if fileErr != nil {
			if errors.As(fileErr, &storage.ErrNoValues) {
				errorMessageResponse(w, "no such file in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, fileErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if editFile.Title == "" {
			editFile.Title = file.Title
		}

		if editFile.File == nil {
			editFile.File = file.File
		}

		if editFile.Notes == "" {
			editFile.Notes = file.Notes
		}

		editFile.ID = file.ID
		editFile.UserID = userID

		newFile, newFileErr := database.Database.EditFile(editFile)
		if newFileErr != nil {
			if errors.As(newFileErr, &storage.ErrNoValues) {
				errorMessageResponse(w, "no such file in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, newFileErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newFile, "application/json", http.StatusCreated)
	}
}

func DeleteFile(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		fileIDStr := chi.URLParam(r, "id")
		fileUUID, fileUUIDErr := uuid.Parse(fileIDStr)
		if fileUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		delErr := database.Database.DeleteFile(fileUUID, userID)
		if delErr != nil {
			if errors.As(delErr, &storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}
			errorMessageResponse(w, delErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, "successful deleted", "application/json", http.StatusOK)
	}
}
