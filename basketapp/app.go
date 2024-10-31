package basketapp

import (
	"git.gocasts.ir/remenu/beehive/basketapp/delivery/http"
	"git.gocasts.ir/remenu/beehive/basketapp/repository"
	"git.gocasts.ir/remenu/beehive/basketapp/service/basket"
	"git.gocasts.ir/remenu/beehive/basketapp/service/order"
	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
	"git.gocasts.ir/remenu/beehive/pkg/logger"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
	"log/slog"
)

type Application struct {
	BasketSvc     basket.Service
	OrderSvc      order.Service
	BasketHandler http.Handler
	BasketRepo    repository.BasketRepo
	HTTPServer    *http.Server
	BasketCfg     Config
	basketLogger  *slog.Logger
}

func Setup(config Config, conn *postgresql.Database) Application {
	// create application struct with all dependencies(repo, broker, delivery)
	// register routes

	return Application{
		HTTPServer:   http.New(*httpserver.New(config.Server)),
		basketLogger: logger.L(),
		BasketCfg:    config,
	}
}

func (app Application) Start() {
	// event listener start
	// long-running process start

	// http/grpc server start

	app.HTTPServer.Serve()

}
