package server

import (
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/handlers"
	"AlexSarva/GophKeeper/internal/app"
	"AlexSarva/GophKeeper/models"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server implementation of custom server
type Server struct {
	httpServer *http.Server
	cfg        *models.Config
	db         *app.Database
}

// NewServer Initializing new server instance
func NewServer() *Server {

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)
	db := *app.NewStorage()
	handler := handlers.CustomHandler(&db)
	server := http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      handler,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
	}
	return &Server{
		httpServer: &server,
		cfg:        &cfg,
		db:         &db,
	}
}

// Run method that starts the server
func (a *Server) Run() error {
	addr := a.httpServer.Addr

	if !a.db.Admin.Ping() {
		log.Fatalln("admin db didnt lunched")
	}
	if !a.db.Repo.Ping() {
		log.Fatalln("work db didnt lunched")
	}

	idleConnsClosed := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-quit
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := a.httpServer.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			log.Printf("Server Shutdown: %v", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	if a.cfg.EnableHTTPS {
		log.Printf("Web-server started at https://%s", addr)
		go func() {
			if err := a.httpServer.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"); err != http.ErrServerClosed {
				log.Fatalf("Failed to listen and serve TLS: %+v", err)
			}
		}()
	} else {
		log.Printf("Web-server started at http://%s", addr)
		go func() {
			if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("Failed to listen and serve: %+v", err)
			}
		}()
	}

	<-idleConnsClosed
	//
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	log.Println("Server Shutdown gracefully")
	return a.httpServer.Shutdown(ctx)
}
