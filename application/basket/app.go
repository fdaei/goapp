package basketapp

import (
	"log/slog"

	baskethttp "git.gocasts.ir/remenu/beehive/delivery/basket/http"
	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
	"git.gocasts.ir/remenu/beehive/pkg/logger"
	"git.gocasts.ir/remenu/beehive/service/basket"
	basketrepo "git.gocasts.ir/remenu/beehive/service/basket/repository"
	"git.gocasts.ir/remenu/beehive/service/order"
)

type Application struct {
	BasketSvc     basket.Service
	OrderSvc      order.Service
	BasketHandler baskethttp.Handler
	BasketRepo    basketrepo.BasketRepo
	HTTPServer    *baskethttp.Server
	BasketCfg     basket.Config
	basketLogger  *slog.Logger
}

func Setup(config basket.Config) Application {
	// create application struct with all dependencies(repo, broker, delivery)
	// register routes

	return Application{
		HTTPServer:   baskethttp.New(*httpserver.New(config.Server)),
		basketLogger: logger.L(),
	}
}

func (app Application) Start() {
	// event listener start
	// long running process start

	// http/grpc server start

	app.HTTPServer.Serve()

}
