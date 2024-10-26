package basket

import (
	netHTTP "net/http"

	"git.gocasts.ir/remenu/beehive/delivery/basket/http"
	basketmysql "git.gocasts.ir/remenu/beehive/repository/mysql/basket"
	"git.gocasts.ir/remenu/beehive/service/basket"
	"git.gocasts.ir/remenu/beehive/service/order"
)

type Application struct {
	BasketSvc     basket.Service
	OrderSvc      order.Service
	BasketHandler http.Handler
	BasketRepo    basketmysql.Repository
	HTTPServer    netHTTP.Server
	BasketCfg     basket.Config
}

func Setup(config basket.Config) Application {
	// create application struct with all dependencies(repo, broker, delivery)
	// register routes
}

func (a Application) Start() {
	// event listener start
	// long running process start
	// http/grpc server start
}
