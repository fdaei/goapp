package basketapp

import (
	"fmt"
	"git.gocasts.ir/remenu/beehive/adapter/rabbitmq"
	"git.gocasts.ir/remenu/beehive/basketapp/delivery/http"
	"git.gocasts.ir/remenu/beehive/basketapp/repository"
	"git.gocasts.ir/remenu/beehive/basketapp/service/basket"
	"git.gocasts.ir/remenu/beehive/basketapp/service/order"
	"git.gocasts.ir/remenu/beehive/event"
	"git.gocasts.ir/remenu/beehive/outbox"
	payment "git.gocasts.ir/remenu/beehive/paymentapp/service"
	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
	"git.gocasts.ir/remenu/beehive/pkg/logger"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
	"log/slog"
	"os"
	"os/signal"
)

type Application struct {
	BasketSvc       basket.Service
	OrderSvc        order.Service
	BasketHandler   http.Handler
	BasketRepo      repository.BasketRepo
	HTTPServer      *http.Server
	BasketCfg       Config
	basketLogger    *slog.Logger
	OutBoxScheduler outbox.Scheduler
}

func Setup(config Config, conn *postgresql.Database) Application {
	// create application struct with all dependencies(repo, broker, delivery)
	// register routes

	outBoxRepo := repository.NewOutBoxRepo(conn.DB)

	queue := "basket"
	topics1 := []event.Topic{
		payment.PurchaseSucceedTopic,
	}
	rabbitMQ := rabbitmq.New(config.RabbitMQ, queue, topics1)

	return Application{
		HTTPServer:      http.New(*httpserver.New(config.Server)),
		basketLogger:    logger.L(),
		BasketCfg:       config,
		OutBoxScheduler: outbox.New(outBoxRepo, rabbitMQ, config.OutboxScheduler),
	}
}

func (app Application) Start() {
	// event listener start
	// long-running process start

	// http/grpc server start

	done := make(chan bool)

	go app.OutBoxScheduler.Start(done)
	go app.HTTPServer.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received interrupt signal, shutting down gracefully..")
	done <- true
}
