package basket

import (
	"log"

	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server struct {
		Port string `koanf:"port"`
	} `koanf:"basket_server"`
	PostgresDB postgresql.Config `koanf:"postgres_basket_db"`
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
