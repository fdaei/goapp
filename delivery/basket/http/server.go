package baskethttp

import (
	"fmt"
	"strconv"

	httpserver "git.gocasts.ir/remenu/beehive/pkg/http_server"
)

type Server struct {
	HTTPServer httpserver.Server
	// TODO: add handlers struct heres in futures
}

func New(server httpserver.Server) *Server {
	return &Server{
		HTTPServer: server,
	}
}

func (s *Server) Serve() {
	s.RegisterRoutes()

	// Start server
	fmt.Printf("start echo server on %s\n", strconv.Itoa(s.HTTPServer.Config.Port))
	if err := s.HTTPServer.Start(); err != nil {
		fmt.Println("router start error", err)
	}
}

func (s *Server) RegisterRoutes() {
	// TODO: add middleware
	// TODO: register swagger

	// Routes
	s.HTTPServer.Router.GET("/health-check", s.healthCheck)
}
