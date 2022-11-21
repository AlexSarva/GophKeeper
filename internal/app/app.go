package app

import (
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/storage"
	"AlexSarva/GophKeeper/storage/admin"
	"AlexSarva/GophKeeper/storage/storagepg"
	"log"
)

// Database interface for different types of databases
type Database struct {
	Repo  storage.Repo
	Admin *admin.Admin
}

// NewStorage generate new instance of database
func NewStorage() *Database {

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	mainStorage := storagepg.PostgresDBConn(cfg.Database)
	adminStorage := admin.NewAdminDBConnection(cfg.AdminDatabase)

	log.Println("Using PostgreSQL Database")

	return &Database{
		Repo:  mainStorage,
		Admin: adminStorage,
	}
}
