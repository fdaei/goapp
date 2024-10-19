package basket

import (
	"log"

	postgresql "git.gocasts.ir/remenu/beehive/pkg/postgresql/config"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	BasketServer struct {
		Port string `koanf:"port"`
	} `koanf:"basket-server"`
	BasketPostgresDB postgresql.Config `koanf:"postgres-basket-db"`
}

func Load(configPath string) *Config {

	var k = koanf.New(".")
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
	return &cfg
}
