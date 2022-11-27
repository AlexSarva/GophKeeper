package handlers

import (
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	var (
		cfg models.ServerConfig
	)
	JSONErr := models.ReadServerJSONConfig(&cfg, "../test/test_server_config.json")
	if JSONErr != nil {
		log.Fatalf("Wrong json format: %+v", JSONErr)
	}
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Fatalln(GlobalContainerErr)
	}
	database := *app.NewStorage()
	tmpUser := models.User{
		ID:       uuid.New(),
		Username: utils.LoginGenerator(7),
		Email:    fmt.Sprintf("%s@gmail.com", utils.LoginGenerator(7)),
		Password: "dPQzaKPD99v",
	}
	user, addUserErr := database.Authorizer.SignUp(tmpUser)
	if addUserErr != nil {
		log.Println(addUserErr)
	}
	user.Password = tmpUser.Password

	type want struct {
		location        string
		contentType     string
		contentEncoding string
		response        string
		code            int
		responseFormat  bool
	}

	tests := []struct {
		name                   string
		request                string
		auth                   string
		protected              bool
		requestPath            string
		requestMethod          string
		requestBody            io.Reader
		requestContentType     string
		requestAcceptEncoding  string
		requestContentEncoding string
		want                   want
	}{
		{
			name:          fmt.Sprintf("%s ping test #1", http.MethodPut),
			requestMethod: http.MethodPut,
			requestPath:   "/ping",
			requestBody:   nil,
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          fmt.Sprintf("%s register negative test #1", http.MethodPost),
			requestMethod: http.MethodPost,
			requestPath:   "/api/v1/register",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"email": "%s@gmail.com",
								"password": "dPQzakp9DMSW"
							}`, utils.LoginGenerator(7)))),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s register negative test #2", http.MethodPost),
			requestMethod: http.MethodPost,
			requestPath:   "/api/v1/register",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"username": "%s",
								"email": "%s@gmail.com",
								"password": "12345"
							}`, utils.LoginGenerator(5), utils.LoginGenerator(5)))),
			want: want{
				code: http.StatusExpectationFailed,
			},
		},
		{
			name:          fmt.Sprintf("%s register positive test #1", http.MethodPost),
			requestMethod: http.MethodPost,
			requestPath:   "/api/v1/register",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"username": "%s",
								"email": "%s@gmail.com",
								"password": "dPQzakp9DMSW"
							}`, utils.LoginGenerator(5), utils.LoginGenerator(5)))),
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:          fmt.Sprintf("%s register negative test #3", http.MethodPost),
			requestMethod: http.MethodPost,
			requestPath:   "/api/v1/register",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"username": "%s",
								"email": "%s",
								"password": "dPQzakp9DMSW"
							}`, utils.LoginGenerator(5), user.Email))),
			want: want{
				code: http.StatusConflict,
			},
		},
		{
			name:          fmt.Sprintf("%s login negative test #1", http.MethodPost),
			requestMethod: http.MethodPost,
			requestPath:   "/api/v1/login",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"email": "%s",
								"password": "2z8PH!1fsaf1"
							}`, user.Email))),
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:          fmt.Sprintf("%s login positive test #1", http.MethodPost),
			requestMethod: http.MethodPost,
			requestPath:   "/api/v1/login",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"email": "%s",
								"password": "%s"
							}`, user.Email, user.Password))),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          fmt.Sprintf("%s user info negative test #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/api/v1/users/me",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"email": "%s",
								"password": "%s"
							}`, user.Email, user.Password))),
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:          fmt.Sprintf("%s user info positive test #1", http.MethodGet),
			requestMethod: http.MethodGet,
			protected:     true,
			auth:          user.Token,
			requestPath:   "/api/v1/users/me",
			requestBody: bytes.NewBuffer([]byte(fmt.Sprintf(`{
								"email": "%s",
								"password": "%s"
							}`, user.Email, user.Password))),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          fmt.Sprintf("%s add cards positive test #1", http.MethodPost),
			requestMethod: http.MethodPost,
			protected:     true,
			auth:          user.Token,
			requestPath:   "/api/v1/info/cards",
			requestBody: bytes.NewBuffer([]byte(`{
    "title": "First card",
    "card_number" : "4405 1111 1000 1383",
    "card_owner" : "Alex sarva",
    "card_exp" : "12/25",
    "notes" : "I go throw all the world"
}`)),
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:          fmt.Sprintf("%s add cards negative test #1", http.MethodPost),
			requestMethod: http.MethodPost,
			protected:     true,
			auth:          user.Token,
			requestPath:   "/api/v1/info/cards",
			requestBody: bytes.NewBuffer([]byte(`{
    "title": "First card",
    "card_owner" : "Alex sarva",
    "card_exp" : "12/25",
    "notes" : "I go throw all the world"
}`)),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          fmt.Sprintf("%s add cards negative test #2", http.MethodPost),
			requestMethod: http.MethodPost,
			requestPath:   "/api/v1/info/cards",
			requestBody: bytes.NewBuffer([]byte(`{
    "title": "First card",
    "card_number" : "4405 1111 1000 1383",
    "card_owner" : "Alex sarva",
    "card_exp" : "12/25",
    "notes" : "I go throw all the world"
}`)),
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}
	//var token string
	Handler := *CustomHandler(&database)
	ts := httptest.NewServer(&Handler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqURL := tt.requestPath + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqURL, tt.requestBody)
			if tt.protected {
				request.Header.Set("Authorization", tt.auth)
			}

			// создаём новый Recorder
			w := httptest.NewRecorder()
			Handler.ServeHTTP(w, request)
			resp := w.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println(err)
				}
			}(resp.Body)

			// Проверяем StatusCode
			respStatusCode := resp.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, wantStatusCode, respStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}
