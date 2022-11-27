package main

import (
	"AlexSarva/GophKeeper/gui"
	"AlexSarva/GophKeeper/models"
	"flag"
	"log"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	cfg          models.GUIConfig
	JSONConfig   models.JSONConfig
)

func version() {
	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)
	log.Println("Warning! Save you id_rsa! if u lost them, u cat'n get any information from database!")
}

func init() {
	flag.StringVar(&cfg.ServerAddress, "address", "", "address of service")
	flag.StringVar(&cfg.KeysPath, "keys", "", "keys filepath")
	flag.IntVar(&cfg.KeysSize, "size", 0, "keys size")
	flag.StringVar(&cfg.Secret, "secret", "", "secret for sym crypt")
	flag.StringVar(&JSONConfig.DSN, "config", "", "JSON config")
}

func main() {
	version()
	flag.Parse()

	if configFilename := JSONConfig.DSN; configFilename != "" {
		JSONErr := models.ReadClientJSONConfig(&cfg, configFilename)
		if JSONErr != nil {
			log.Fatalf("Wrong json format: %+v", JSONErr)
		}
	}

	if cfg.ServerAddress == "" {
		log.Fatalln("cant obtain server address")
	}

	if cfg.Secret == "" {
		log.Fatalln("cant obtain secret for sym crypto")
	}

	myGUI := gui.InitGUI(&cfg)
	//gui.Render()
	myGUI.Render()
	runErr := myGUI.Run()
	if runErr != nil {
		log.Fatalln(runErr)
	}
}
