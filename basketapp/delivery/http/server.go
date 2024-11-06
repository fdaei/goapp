package http

import (
	"fmt"
	"strconv"

	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
)

type Server struct {
	HTTPServer httpserver.Server
	Handler    Handler
}

func New(server httpserver.Server, handler Handler) Server {
	return Server{
		HTTPServer: server,
		Handler:    handler,
	}
}

func (s Server) Serve() {
	s.RegisterRoutes()

	// Start server
	fmt.Printf("start echo server on %s\n", strconv.Itoa(s.HTTPServer.Config.Port))
	if err := s.HTTPServer.Start(); err != nil {
		fmt.Println("router start error", err)
	}
}

func (s Server) RegisterRoutes() {
	// TODO: add middleware
	// TODO: register swagger

	// Define Routes

	//health check
	s.HTTPServer.Router.GET("/health-check", s.healthCheck)

	// basket group
	basketGroup := s.HTTPServer.Router.Group("/basket")
	// basket CRUD
	basketGroup.POST("/add", s.Handler.AddToBasket)
}
