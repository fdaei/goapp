package config

import "git.gocasts.ir/remenu/beehive/adapter/rabbitmq"

type Config struct {
	RabbitMQ rabbitmq.Config `koanf:"rabbitmq"`
	BasketDB BasketDBConfig  `koanf:"postgres-basket-db"`
}

type BasketDBConfig struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	User            string `koanf:"user"`
	Password        string `koanf:"password"`
	DBName          string `koanf:"dbName"`
	SSLMode         string   `koanf:"sslMode"`
	MaxIdleConns    int    `koanf:"maxIdleConns"`
	MaxOpenConns    int    `koanf:"maxOpenConns"`
	ConnMaxLifetime int    `koanf:"connMaxLifetime"`
}
