package config

import "git.gocasts.ir/remenu/beehive/adapter/rabbitmq"

type Config struct {
	RabbitMQ rabbitmq.Config `koanf:"rabbitmq"`
}
