package handlers

import (
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRegistration - user registration method
//
// Handler POST /api/v1/admin/register
//
// Registration is performed by a pair of login/password.
// Each login must be set.
// After successful registration, automatic user authentication is required.
// post message should contain such body:
//
//	"login": "<login>",
//	"email": "<email>",
//	"password": "<password>"
//
// Possible response codes:
// 201 - user successfully registered and authenticated;
// 400 - invalid request format;
// 409 - login is already taken;
// 500 - an internal server error.
func UserRegistration(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		readBodyErr := readBodyInStruct(r, &user)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}
		if user.Email == "" || user.Password == "" || user.Username == "" {
			errorMessageResponse(w, "empty fields error", "application/json", http.StatusBadRequest)
			return
		}

		userID := uuid.New()
		userToken, userTokenExp := GenerateToken(userID)
		hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
		if bcrypteErr != nil {
			errorMessageResponse(w, ErrCryptPassword.Error(), "application/json", http.StatusBadRequest)
			return
		}

		user.ID, user.Password, user.Token, user.TokenExp = userID, string(hashedPassword), userToken, userTokenExp

		newUserErr := database.Admin.Register(user)
		if newUserErr != nil {
			if newUserErr == storage.ErrDuplicatePK {
				errorMessageResponse(w, ErrLoginExist.Error(), "application/json", http.StatusConflict)
				return
			}
			errorMessageResponse(w, newUserErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		userInfo, userInfoErr := database.Admin.GetUserInfo(userID)
		if userInfoErr != nil {
			errorMessageResponse(w, userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, userInfo, "application/json", http.StatusCreated)
	}
}

// UserAuthentication - user authentication method
//
// Handler POST /api/v1/admin/login
//
// Authentication is performed by a login/password pair.
// Request format:
//
//	{"login": "<login>",
//	"password": "<password>"}
//
// Possible response codes:
// 200 - user successfully authenticated;
// 400 - invalid request format;
// 401 - invalid login/password pair;
// 500 - an internal server error.
func UserAuthentication(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.UserLogin
		readBodyErr := readBodyInStruct(r, &user)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}

		passwdDB, passwdDBErr := database.Admin.Login(user.Email)
		if passwdDBErr != nil {
			if errors.Is(passwdDBErr, sql.ErrNoRows) {
				errorMessageResponse(w, ErrNoUserExists.Error(), "application/json", http.StatusUnauthorized)
				return
			}
			errorMessageResponse(w, passwdDBErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		cryptErr := bcrypt.CompareHashAndPassword([]byte(passwdDB.Password), []byte(user.Password))
		if cryptErr != nil {
			errorMessageResponse(w, ErrComparePassword.Error(), "application/json", http.StatusUnauthorized)
			return
		}
		// TODO Предусмотреть обновление куки

		userInfo, userInfoErr := database.Admin.GetUserInfo(passwdDB.ID)
		if userInfoErr != nil {
			errorMessageResponse(w, "Internal Server Error "+userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, userInfo, "application/json", http.StatusOK)
	}
}
