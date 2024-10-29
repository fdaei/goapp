package main

import (
	"log"
	"os"
	"path/filepath"

	"git.gocasts.ir/remenu/beehive/basketapp"
	cfgloader "git.gocasts.ir/remenu/beehive/pkg/cfg_loader"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
)

var cfg basketapp.Config

func main() {
	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	options := cfgloader.Option{
		Prefix:       "BASKET_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "basket", "development", "config.yaml"),
		CallbackEnv:  nil,
	}

	if err := cfgloader.Load(options, &cfg); err != nil {
		log.Fatalf("Failed to load basket config: %v", err)
	}

	// show loaded config
	log.Printf("Loaded config: %+v\n", cfg)

	// Connect to database
	conn, cnErr := postgresql.Connect(cfg.PostgresDB)

	if cnErr != nil {
		log.Fatal(cnErr)
	} else {
		log.Printf("You are connected to %s successfully\n", cfg.PostgresDB.DBName)
	}

	// Close the database connection
	defer postgresql.Close(conn.DB)

	// start application
	app := basketapp.Setup(cfg, conn)
	app.Start()

}
