package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/handler/rest/client"
	"github.com/tuanta7/ciam/internal/handler/rest/login"
	"github.com/tuanta7/ciam/internal/handler/rest/middleware"
	"github.com/tuanta7/ciam/internal/usecase/oidc"
)

type Server struct {
	server        *http.Server
	router        chi.Router
	op            *oidc.Provider
	authenticator *login.Handler
	clientHandler *client.Handler
}

func NewServer(
	cfg *config.EnvConfig,
	provider *oidc.Provider,
	authenticator *login.Handler,
	clientHandler *client.Handler,
) *Server {
	router := chi.NewRouter()

	return &Server{
		server: &http.Server{
			Addr:    cfg.BindAddress,
			Handler: router,
		},
		router:        router,
		op:            provider,
		authenticator: authenticator,
		clientHandler: clientHandler,
	}
}

func (s *Server) Run() error {
	s.router.Mount("/", s.op)
	s.router.Get("/login", s.authenticator.GetLoginForm)
	s.router.Post("/login", s.authenticator.Login)

	s.router.Route("/api/v1/clients", func(r chi.Router) {
		r.With(middleware.Pagination).Get("/", s.clientHandler.ListClients)
		r.Post("/", s.clientHandler.CreateClient)
		r.Get("/{id}", s.clientHandler.GetClient)
		r.Put("/{id}", s.clientHandler.UpdateClient)
		r.Delete("/{id}", s.clientHandler.DeleteClient)
	})

	log.Printf("OP server is running on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Timeout() time.Duration {
	return 20 * time.Second
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
