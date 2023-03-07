package app

import (
	"AlexSarva/GophKeeper/authorizer"
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"AlexSarva/GophKeeper/storage/admin"
	"AlexSarva/GophKeeper/storage/storagepg"
	"AlexSarva/GophKeeper/utils"
	"log"
	"time"
)

// Storage interface for different types of databases
type Storage struct {
	Database        storage.Database
	Admin           *admin.Admin
	Authorizer      *authorizer.Authorizer
	PasswordChecker *utils.PasswordChecker
}

// NewStorage generate new instance of database
func NewStorage() *Storage {

	cfg := constant.GlobalContainer.Get("server-config").(models.ServerConfig)

	mainStorage := storagepg.PostgresDBConn(cfg.Database)
	adminStorage := admin.NewAdminDBConnection(cfg.AdminDatabase)
	auth := authorizer.NewAuthorizer(adminStorage, []byte(cfg.Secret), 72*time.Hour)
	passwordChecker := utils.InitPasswordChecker(8, true, true, false)

	log.Println("Using PostgreSQL Database")

	return &Storage{
		Database:        mainStorage,
		Admin:           adminStorage,
		Authorizer:      auth,
		PasswordChecker: passwordChecker,
	}
}
