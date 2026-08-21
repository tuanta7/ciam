package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/handler/rest/client"
	"github.com/tuanta7/ciam/internal/handler/rest/middleware"
)

type Server struct {
	server        *http.Server
	router        chi.Router
	clientHandler *client.Handler
}

func NewServer(
	cfg *config.EnvConfig,
	clientHandler *client.Handler,
) *Server {
	router := chi.NewRouter()

	return &Server{
		server: &http.Server{
			Addr:    cfg.AdminBindAddress,
			Handler: router,
		},
		router:        router,
		clientHandler: clientHandler,
	}
}

func (s *Server) Run() error {
	s.router.Route("/api/v1/clients", func(r chi.Router) {
		r.With(middleware.Pagination).Get("/", s.clientHandler.ListClients)
		r.Post("/", s.clientHandler.CreateClient)
		r.Get("/{id}", s.clientHandler.GetClient)
		r.Put("/{id}", s.clientHandler.UpdateClient)
		r.Delete("/{id}", s.clientHandler.DeleteClient)
	})

	log.Printf("Admin server is running on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Timeout() time.Duration {
	return 20 * time.Second
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
