package main

import (
	"AlexSarva/GophKeeper/constant"
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/server"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	cfg          models.Config
	JSONConfig   models.JSONConfig
)

func version() {
	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)
}

func init() {
	flag.StringVar(&cfg.ServerAddress, "address", "", "host:port to listen on")
	flag.StringVar(&cfg.Database, "database", "", "database config")
	flag.StringVar(&cfg.AdminDatabase, "admin", "", "admin database config")
	flag.StringVar(&cfg.Secret, "secret", "", "secret word")
	flag.StringVar(&cfg.CORS, "cors", "", "cors settings")
	flag.StringVar(&JSONConfig.DSN, "config", "", "JSON config")
	flag.BoolVar(&cfg.EnableHTTPS, "secure", false, "enable HTTPS")
	flag.StringVar(&cfg.TrustedSubnet, "trusted", "", "trusted subnet")
}

func main() {
	version()

	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Перезаписываем из параметров запуска
	flag.Parse()

	if cfg.Database == "" && !cfg.EnableHTTPS && cfg.TrustedSubnet == "" && cfg.Secret == "" && cfg.CORS == "" {
		if configFilename := JSONConfig.DSN; configFilename != "" {
			JSONErr := models.ReadJSONConfig(&cfg, configFilename)
			if JSONErr != nil {
				log.Fatalf("Wrong json format: %+v", JSONErr)
			}
		}
	}

	log.Printf("ServerAddress: %v, EnableHTTPS: %v", cfg.ServerAddress, cfg.EnableHTTPS)

	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Fatalln(GlobalContainerErr)
	}

	MainApp := server.NewServer()
	if errApp := MainApp.Run(); errApp != nil {
		log.Fatalln(errApp)
	}
}
