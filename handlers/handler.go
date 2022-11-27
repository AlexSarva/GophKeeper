package handlers

import (
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	ErrJSONWrite    = errors.New("problem in writing json")
	ErrJSONRequest  = errors.New("wrong type provided for fields")
	ErrLoginExist   = errors.New("login is already busy")
	ErrUnauthorized = errors.New("user unauthorized")
)

// errorMessageResponse additional respond generator
// useful in case of error handling in outputting results to respond
func errorMessageResponse(w http.ResponseWriter, message, ContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(httpStatusCode)
	resp := map[string]string{"message": message}
	jsonResp, jsonRespErr := json.Marshal(resp)
	if jsonRespErr != nil {
		log.Println(jsonRespErr)
	}
	_, writeErr := w.Write(jsonResp)
	if writeErr != nil {
		log.Println("something wrong happens", writeErr)
	}
}

// resultResponse additional result response generator
func resultResponse(w http.ResponseWriter, data interface{}, ContentType string, httpStatusCode int) {
	jsonResp, jsonRespErr := json.Marshal(data)
	if jsonRespErr != nil {
		errorMessageResponse(w, ErrJSONWrite.Error(), "application/json", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(httpStatusCode)
	_, writeErr := w.Write(jsonResp)
	if writeErr != nil {
		log.Println("something wrong happens", writeErr)
	}
}

// readBodyInStruct read compressed and usual request in struct
func readBodyInStruct(r *http.Request, data interface{}) error {
	// GZIP decode
	var body io.ReadCloser
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			return readErr
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(r.Body)

		newR, gzErr := gzip.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return gzErr
		}
		defer func(newR *gzip.Reader) {
			err := newR.Close()
			if err != nil {
				log.Println(err)
			}
		}(newR)

		body = newR
	} else {
		body = r.Body
	}

	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	errDecode := decoder.Decode(&data)

	if errDecode != nil {
		if errors.Is(errDecode, unmarshalErr) {
			return ErrJSONRequest
		}
		return errDecode
	}
	return nil
}

// gzipContentTypes request types that support data compression
var gzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

func CustomAllowOriginFunc(_ *http.Request, origin string) bool {
	cfg := constant.GlobalContainer.Get("server-config").(models.ServerConfig)
	urls := strings.Fields(cfg.CORS)
	corsMap := make(map[string]bool)
	for i := 0; i < len(urls); i += 1 {
		corsMap[urls[i]] = true
	}
	return corsMap[origin]
}

// Ping returns pong
func Ping(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Println(err)
	}
}

// CustomHandler - the main_admin_test handler of the server
// contains middlewares and all routes
func CustomHandler(database *app.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: CustomAllowOriginFunc,
		//AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.AllowContentEncoding("gzip"))
	//r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))
	r.Use(checkContent)
	r.Mount("/debug", middleware.Profiler())
	//
	r.Put("/ping", Ping)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", UserRegistration(database))
		r.Post("/login", UserAuthentication(database))

		r.Route("/users", func(r chi.Router) {
			r.Use(userIdentification(database))
			r.Get("/me", GetUserInfo(database))
		})

		r.Route("/info", func(r chi.Router) {
			r.Use(userIdentification(database))
			r.Route("/notes", func(r chi.Router) {
				r.Get("/", GetNoteList(database))
				r.Post("/", PostNote(database))
				r.Get("/{id}", GetNote(database))
				r.Patch("/{id}", EditNote(database))
				r.Delete("/{id}", DeleteNote(database))
			})
			r.Route("/cards", func(r chi.Router) {
				r.Get("/", GetCardList(database))
				r.Post("/", PostCard(database))
				r.Get("/{id}", GetCard(database))
				r.Patch("/{id}", EditCard(database))
				r.Delete("/{id}", DeleteCard(database))
			})
			r.Route("/creds", func(r chi.Router) {
				r.Get("/", GetCredList(database))
				r.Post("/", PostCred(database))
				r.Get("/{id}", GetCred(database))
				r.Patch("/{id}", EditCred(database))
				r.Delete("/{id}", DeleteCred(database))
			})
			r.Route("/files", func(r chi.Router) {
				r.Get("/", GetFileList(database))
				r.Post("/", PostFile(database))
				r.Get("/{id}", GetFile(database))
				r.Patch("/{id}", EditFile(database))
				r.Delete("/{id}", DeleteFile(database))
			})
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, this method are not allowed."))
		if naErr != nil {
			log.Println(naErr)
		}
	})
	return r
}
