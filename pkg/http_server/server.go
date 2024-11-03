package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port int  `koanf:"port"`
	Cors Cors `koanf:"cors"`
}

type Cors struct {
	AllowOrigins []string `koanf:"allow_origins"`
}

type Server struct {
	Router *echo.Echo
	Config Config
}

func New(cfg Config) *Server {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Cors.AllowOrigins,
	}))

	return &Server{
		Router: e,
		Config: cfg,
	}
}

// register custom handler
func (s *Server) RegisterHandler(route string, handler echo.HandlerFunc) {
	s.Router.GET(route, handler)
}

// start server
func (s *Server) Start() error {
	return s.Router.Start(fmt.Sprintf(":%d", s.Config.Port))
}