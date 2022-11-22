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

func PostNote(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var note models.NewNote
		readBodyErr := readBodyInStruct(r, &note)
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
		note.UserID = userID

		newNote, newNoteErr := database.Database.NewNote(&note)
		if newNoteErr != nil {
			errorMessageResponse(w, newNoteErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		resultResponse(w, newNote, "application/json", http.StatusCreated)
	}
}

func GetNoteList(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		notes, notesErr := database.Database.AllNotes(userID)
		if notesErr != nil {
			errorMessageResponse(w, notesErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		if len(notes) == 0 {
			errorMessageResponse(w, "no values", "application/json", http.StatusNoContent)
			return
		}

		resultResponse(w, notes, "application/json", http.StatusOK)
	}
}

func GetNote(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		noteIDStr := chi.URLParam(r, "id")
		noteUUID, noteUUIDErr := uuid.Parse(noteIDStr)
		if noteUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		note, notesErr := database.Database.GetNote(noteUUID, userID)
		if notesErr != nil {
			if errors.As(notesErr, &storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, notesErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		resultResponse(w, note, "application/json", http.StatusOK)
	}
}

func EditNote(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editNote models.NewNote
		readBodyErr := readBodyInStruct(r, &editNote)
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

		noteIDStr := chi.URLParam(r, "id")
		noteUUID, noteUUIDErr := uuid.Parse(noteIDStr)
		if noteUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		note, notesErr := database.Database.GetNote(noteUUID, userID)
		if notesErr != nil {
			if errors.As(notesErr, &storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, notesErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if editNote.Title == "" {
			editNote.Title = note.Title
		}

		if editNote.Note == "" {
			editNote.Note = note.Note
		}

		editNote.ID = note.ID
		editNote.UserID = userID

		newNote, newNoteErr := database.Database.EditNote(editNote)
		if newNoteErr != nil {
			if errors.As(newNoteErr, &storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, newNoteErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newNote, "application/json", http.StatusCreated)
	}
}

func DeleteNote(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		noteIDStr := chi.URLParam(r, "id")
		noteUUID, noteUUIDErr := uuid.Parse(noteIDStr)
		if noteUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		delErr := database.Database.DeleteNote(noteUUID, userID)
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
