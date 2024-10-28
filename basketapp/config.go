package basketapp

import (
	"git.gocasts.ir/remenu/beehive/basketapp/service/basket"
	"git.gocasts.ir/remenu/beehive/basketapp/service/order"
	"git.gocasts.ir/remenu/beehive/pkg/logger"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
)

type Config struct {
	BasketSvcCfg basket.Config
	OrderSvcCfg  order.Config
	Server       struct {
		Port string `koanf:"port"`
	} `koanf:"server"`
	PostgresDB postgresql.Config `koanf:"basket_postgres_db"`
	Logger     logger.Config     `koanf:"logger"`
}
