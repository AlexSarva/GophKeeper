package models

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// Config  start parameters for lunch the service
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	Database      string `env:"DATABASE_DSN" json:"database_dsn"`
	AdminDatabase string `env:"ADMIN_DATABASE_DSN" json:"admin_database_dsn"`
	Secret        string `env:"SECRET" json:"secret"`
	CORS          string `env:"CORS" json:"cors"`
	EnableHTTPS   bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// JSONConfig config file in json format
type JSONConfig struct {
	DSN string
}

func ReadJSONConfig(cfg *Config, JSONFilepath string) error {
	f, fErr := os.Open(JSONFilepath)
	log.Println("read lunch parameters from cfg file")
	if fErr != nil {
		return fErr
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)

	var unmarshalConfigErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&cfg)

	if err != nil {
		if errors.As(err, &unmarshalConfigErr) {
			return err
		} else {
			return err
		}
	}

	return nil
}
