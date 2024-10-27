package basket

import (
	"log/slog"
	netHTTP "net/http"

	"git.gocasts.ir/remenu/beehive/delivery/basket/http"
	"git.gocasts.ir/remenu/beehive/service/basket"
	basketrepo "git.gocasts.ir/remenu/beehive/service/basket/repository"
	"git.gocasts.ir/remenu/beehive/service/order"
)

type Application struct {
	BasketSvc     basket.Service
	OrderSvc      order.Service
	BasketHandler http.Handler
	BasketRepo    basketrepo.BasketRepo
	HTTPServer    netHTTP.Server
	BasketCfg     basket.Config
	basketLogger  *slog.Logger
}

func Setup(config basket.Config) Application {
	// create application struct with all dependencies(repo, broker, delivery)
	// register routes
	return Application{}
}

func (a Application) Start() {
	// event listener start
	// long running process start
	// http/grpc server start
}
