package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cfgloader "git.gocasts.ir/remenu/beehive/pkg/cfg_loader"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
	"git.gocasts.ir/remenu/beehive/service/basket"
)

var cfg basket.Config

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
	fmt.Printf("Loaded config: %+v\n", cfg)

	// Connect to database
	conn, cnErr := postgresql.Connect(cfg.PostgresDB)

	if cnErr != nil {
		log.Fatal(cnErr)
	} else {
		fmt.Printf("You are connected to %s successfully\n", cfg.PostgresDB.DBName)
	}

	// Check example query to ensure that db works correctly
	res, exErr := postgresql.ExampleQuery(conn.DB)
	if exErr != nil {
		log.Fatal(exErr)
	} else {
		fmt.Printf("The version of database is: %s\n", res)
	}

	// Close the database connection
	defer postgresql.Close(conn.DB)

}
