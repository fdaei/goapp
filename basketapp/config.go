package basketapp

import (
	"git.gocasts.ir/remenu/beehive/adapter/rabbitmq"
	"git.gocasts.ir/remenu/beehive/basketapp/service/basket"
	"git.gocasts.ir/remenu/beehive/basketapp/service/order"
	"git.gocasts.ir/remenu/beehive/outbox"
	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
	"git.gocasts.ir/remenu/beehive/pkg/logger"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
)

type Config struct {
	BasketSvcCfg    basket.Config
	OrderSvcCfg     order.Config
	Server          httpserver.Config `koanf:"server"`
	PostgresDB      postgresql.Config `koanf:"postgres_db"`
	Logger          logger.Config     `koanf:"logger"`
	OutboxScheduler outbox.Config     `koanf:"outbox_scheduler"`
	RabbitMQ        rabbitmq.Config   `koanf:"rabbitmq"`
}
