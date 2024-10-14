package main

import (
	"fmt"
	"log"

	"git.gocasts.ir/remenu/beehive/config"
	basketpostgres "git.gocasts.ir/remenu/beehive/repository/postgresql/basket"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func main() {

	if err := k.Load(file.Provider("config/basket.yml"), yaml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var cfg config.Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	// Connect to database
	conn, cnErr := basketpostgres.Connect(cfg.BasketDB)

	if cnErr != nil {
		log.Fatal(cnErr)
	} else {
		fmt.Printf("You are connected to %s successfully\n", cfg.BasketDB.DBName)
	}

	// Check example query to ensure that db works correctly
	res, exErr := basketpostgres.ExampleQuery(conn.DB)
	if exErr != nil {
		log.Fatal(exErr)
	} else {
		fmt.Printf("The version of database is: %s\n", res)
	}

	// Close the database connection
	defer basketpostgres.Close(conn.DB)

}
