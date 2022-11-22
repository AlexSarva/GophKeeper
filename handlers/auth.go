package handlers

import (
	"AlexSarva/GophKeeper/authorizer"
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"AlexSarva/GophKeeper/utils"
	"errors"
	"net/http"

	"github.com/google/uuid"
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
func UserRegistration(database *app.Storage) http.HandlerFunc {
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

		user.ID = uuid.New()
		newUser, newUserErr := database.Authorizer.SignUp(user)

		if newUserErr != nil {
			if errors.As(newUserErr, &storage.ErrDuplicatePK) {
				errorMessageResponse(w, ErrLoginExist.Error(), "application/json", http.StatusConflict)
				return
			}

			if errors.As(newUserErr, &authorizer.ErrHashPassword) || errors.As(newUserErr, &authorizer.ErrGenerateToken) {
				errorMessageResponse(w, newUserErr.Error(), "application/json", http.StatusBadRequest)
				return
			}

			errorMessageResponse(w, newUserErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, newUser, "application/json", http.StatusCreated)
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
func UserAuthentication(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.UserLogin
		readBodyErr := readBodyInStruct(r, &user)
		if readBodyErr != nil {
			errorMessageResponse(w, readBodyErr.Error(), "application/json", http.StatusBadRequest)
			return
		}

		userInfo, userInfoErr := database.Authorizer.SignIn(&user)
		if userInfoErr != nil {
			if errors.As(userInfoErr, &authorizer.ErrNoUserExists) {
				errorMessageResponse(w, userInfoErr.Error(), "application/json", http.StatusUnauthorized)
				return
			}
			if errors.As(userInfoErr, &authorizer.ErrComparePassword) {
				errorMessageResponse(w, userInfoErr.Error(), "application/json", http.StatusUnauthorized)
				return
			}
			errorMessageResponse(w, userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, userInfo, "application/json", http.StatusOK)
	}
}

// GetUserInfo - user info method
//
// Handler GET /api/v1/admin/users/me
//
// Authorization: "Bearer T"
//
// Possible response codes:
// 200 - load user information;
// 400 - invalid request format;
// 401 - invalid auth;
// 500 - an internal server error.
func GetUserInfo(database *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userIDErr := utils.GetUserID(ctx)
		if userIDErr != nil {
			errorMessageResponse(w, ErrUnauthorized.Error()+": "+userIDErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		userInfo, userInfoErr := database.Admin.GetUserInfo(userID)
		if userInfoErr != nil {
			errorMessageResponse(w, userInfoErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		resultResponse(w, userInfo, "application/json", http.StatusOK)
	}
}
