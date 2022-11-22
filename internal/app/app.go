package app

import (
	"AlexSarva/GophKeeper/authorizer"
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"AlexSarva/GophKeeper/storage/admin"
	"AlexSarva/GophKeeper/storage/storagepg"
	"log"
	"time"
)

// Storage interface for different types of databases
type Storage struct {
	Database   storage.Database
	Admin      *admin.Admin
	Authorizer *authorizer.Authorizer
}

// NewStorage generate new instance of database
func NewStorage() *Storage {

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	mainStorage := storagepg.PostgresDBConn(cfg.Database)
	adminStorage := admin.NewAdminDBConnection(cfg.AdminDatabase)
	auth := authorizer.NewAuthorizer(adminStorage, []byte(cfg.Secret), 72*time.Hour)

	log.Println("Using PostgreSQL Database")

	return &Storage{
		Database:   mainStorage,
		Admin:      adminStorage,
		Authorizer: auth,
	}
}
