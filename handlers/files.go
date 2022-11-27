package handlers

import (
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func PostFile(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		if title == "" {
			errorMessageResponse(w, "dont have parameter 'title' in request", "application/json", http.StatusBadRequest)
			return
		}
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			errorMessageResponse(w, "dont have parameter 'filename' in request", "application/json", http.StatusBadRequest)
			return
		}
		notes := r.URL.Query().Get("notes")
		var file models.NewFile
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				errorMessageResponse(w, err.Error(), "application/json", http.StatusBadRequest)
				return
			}
		}(r.Body)
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("request", err)
		}
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}
		file.File = buf
		file.Title = title
		file.FileName = filename
		if notes != "" {
			file.Notes = notes
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
		userID, userIDErr := GetUserID(ctx)
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
		userID, userIDErr := GetUserID(ctx)
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
			if errors.Is(fileErr, storage.ErrNoValues) {
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
		title := r.URL.Query().Get("title")
		if title == "" {
			errorMessageResponse(w, "dont have parameter 'title' in request", "application/json", http.StatusBadRequest)
			return
		}
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			errorMessageResponse(w, "dont have parameter 'filename' in request", "application/json", http.StatusBadRequest)
			return
		}
		notes := r.URL.Query().Get("notes")
		var editFile models.NewFile
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				errorMessageResponse(w, err.Error(), "application/json", http.StatusBadRequest)
				return
			}
		}(r.Body)
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("request", err)
		}
		editFile.File = buf
		editFile.Title = title
		editFile.FileName = filename
		if notes != "" {
			editFile.Notes = notes
		}
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
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
			if errors.Is(fileErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such file in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, fileErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if editFile.File == nil {
			editFile.File = file.File
		}

		if editFile.Notes == "" {
			editFile.Notes = file.Notes
		}

		editFile.ID = file.ID
		editFile.UserID = userID

		newFile, newFileErr := database.Database.EditFile(&editFile)
		if newFileErr != nil {
			if errors.Is(newFileErr, storage.ErrNoValues) {
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
		userID, userIDErr := GetUserID(ctx)
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
