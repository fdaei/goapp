package main

import (
	"fmt"
	"git.gocasts.ir/remenu/beehive/outbox/repository"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"git.gocasts.ir/remenu/beehive/adapter/rabbitmq"
	"git.gocasts.ir/remenu/beehive/basketapp"
	"git.gocasts.ir/remenu/beehive/event"
	"git.gocasts.ir/remenu/beehive/outbox"
	payment "git.gocasts.ir/remenu/beehive/paymentapp/service"
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

	outBoxRepo := repository.NewOutBoxRepo(conn.DB)

	queue := "basket"
	topics1 := []event.Topic{
		payment.PurchaseSucceedTopic,
	}
	rabbitMQ := rabbitmq.New(cfg.RabbitMQ, queue, topics1)

	done := make(chan bool)

	sch := outbox.New(outBoxRepo, rabbitMQ, cfg.OutboxScheduler)
	sch.Start(done)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received interrupt signal, shutting down gracefully..")
	done <- true

}
