package basket

import (
	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
	"git.gocasts.ir/remenu/beehive/pkg/logger"
	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
)

type Config struct {
	Server     httpserver.Config `koanf:"server"`
	PostgresDB postgresql.Config `koanf:"postgres_db"`
	Logger     logger.Config     `koanf:"logger"`
}
