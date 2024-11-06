package basketapp

import (
	"context"
	"log/slog"

	"git.gocasts.ir/remenu/beehive/basketapp/delivery/http"
	"git.gocasts.ir/remenu/beehive/basketapp/repository"
	"git.gocasts.ir/remenu/beehive/basketapp/service/basket"
	"git.gocasts.ir/remenu/beehive/basketapp/service/order"
	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
	"git.gocasts.ir/remenu/beehive/pkg/logger"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
)

type Application struct {
	Ctx           context.Context
	BasketRepo    basket.Repository
	BasketSvc     basket.Service
	BasketHandler http.Handler
	OrderSvc      order.Service
	HTTPServer    http.Server
	BasketCfg     Config
	BasketLogger  *slog.Logger
}

func Setup(ctx context.Context, config Config, conn *postgresql.Database) Application {
	// create application struct with all dependencies(repo, broker, delivery)
	basketRepo := repository.NewBasketRepo(conn.DB)
	basketSvc := basket.NewService(basketRepo)
	basketHandler := http.NewHandler(basketSvc)

	return Application{
		Ctx:          ctx,
		BasketRepo:   basketRepo,
		BasketSvc:    basketSvc,
		BasketHandler: basketHandler,
		HTTPServer:   http.New(httpserver.New(config.Server), basketHandler),
		BasketLogger: logger.L(),
		BasketCfg:    config,
	}
}

func (app Application) Start() {
	// event listener start
	// long-running process start

	// http/grpc server start

	app.HTTPServer.Serve()

}
