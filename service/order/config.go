package order

import (
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
)

type Config struct {
	Server struct {
		Port string `koanf:"port"`
	} `koanf:"server"`
	PostgresDB postgresql.Config `koanf:"postgres_db"`
}
