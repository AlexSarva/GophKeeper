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

func PostCard(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var card models.NewCard
		readBodyErr := readBodyInStruct(r, &card)
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
		card.UserID = userID

		newNote, newNoteErr := database.Database.NewCard(&card)
		if newNoteErr != nil {
			errorMessageResponse(w, newNoteErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newNote, "application/json", http.StatusCreated)
	}
}

func GetCardList(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		cards, notesErr := database.Database.AllCards(userID)
		if notesErr != nil {
			errorMessageResponse(w, notesErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		if len(cards) == 0 {
			errorMessageResponse(w, "no values", "application/json", http.StatusNoContent)
			return
		}

		resultResponse(w, cards, "application/json", http.StatusOK)
	}
}

func GetCard(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		cardIDStr := chi.URLParam(r, "id")
		cardUUID, cardUUIDErr := uuid.Parse(cardIDStr)
		if cardUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		card, notesErr := database.Database.GetCard(cardUUID, userID)
		if notesErr != nil {
			if errors.Is(notesErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such note in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, notesErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}
		resultResponse(w, card, "application/json", http.StatusOK)
	}
}

func EditCard(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editCard models.NewCard
		readBodyErr := readBodyInStruct(r, &editCard)
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

		cardIDStr := chi.URLParam(r, "id")
		cardUUID, cardUUIDErr := uuid.Parse(cardIDStr)
		if cardUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		card, notesErr := database.Database.GetCard(cardUUID, userID)
		if notesErr != nil {
			if errors.Is(notesErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such card in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, notesErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if editCard.Title == "" {
			editCard.Title = card.Title
		}

		if editCard.CardNumber == "" {
			editCard.CardNumber = card.CardNumber
		}

		if editCard.CardOwner == "" {
			editCard.CardOwner = card.CardOwner
		}

		if editCard.CardExp == "" {
			editCard.CardExp = card.CardExp
		}

		if editCard.Notes == "" {
			editCard.Notes = card.Notes
		}

		editCard.ID = card.ID
		editCard.UserID = userID

		newCard, newCardErr := database.Database.EditCard(editCard)
		if newCardErr != nil {
			if errors.Is(newCardErr, storage.ErrNoValues) {
				errorMessageResponse(w, "no such card in db", "application/json", http.StatusConflict)
				return
			}

			errorMessageResponse(w, newCardErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newCard, "application/json", http.StatusCreated)
	}
}

func DeleteCard(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		cardIDStr := chi.URLParam(r, "id")
		cardUUID, cardUUIDErr := uuid.Parse(cardIDStr)
		if cardUUIDErr != nil {
			errorMessageResponse(w, "Check ID please", "application/json", http.StatusBadRequest)
			return
		}

		delErr := database.Database.DeleteCard(cardUUID, userID)
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
