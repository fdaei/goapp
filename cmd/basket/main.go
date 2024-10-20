package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
	"git.gocasts.ir/remenu/beehive/service/basket"
)

func main() {
	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	//configPath := filepath.Join(workingDir, "deploy", "basket", "development", "config.yaml")
	configPath := GetConfigPath(workingDir, "basket", "development")
	basketCfg := basket.Load(configPath)
	// show loaded config
	fmt.Printf("Loaded config: %+v\n", basketCfg)

	// Connect to database
	conn, cnErr := postgresql.Connect(basketCfg.PostgresDB)

	if cnErr != nil {
		log.Fatal(cnErr)
	} else {
		fmt.Printf("You are connected to %s successfully\n", basketCfg.PostgresDB.DBName)
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

func GetConfigPath(workingDir, serviceName, environment string) string {

	return filepath.Join(workingDir, "deploy", serviceName, environment, "config.yaml")
}
