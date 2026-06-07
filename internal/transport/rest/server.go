package rest

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tuanta7/ciam/internal/transport/rest/handler"
	"github.com/tuanta7/ciam/internal/transport/rest/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Server struct {
	server        *http.Server
	router        chi.Router
	meter         metric.Meter
	clientHandler *handler.ClientHandler
}

func NewServer(
	addr string,
	clientHandler *handler.ClientHandler,
) *Server {
	router := chi.NewRouter()

	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		router:        router,
		clientHandler: clientHandler,
		meter:         otel.Meter("rest_server_meter"),
	}
}

func (s *Server) Run() error {
	if err := middleware.InitMetricsMiddleware(s.meter); err != nil {
		return err
	}

	s.registerRoutes()

	log.Printf("Server is running on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Timeout() time.Duration {
	return 20 * time.Second
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	s.router.Route("/api/internal/v1/clients", func(r chi.Router) {
		r.Use(middleware.WithMetric)

		r.With(middleware.Pagination).Get("/", s.clientHandler.ListClients)
		r.Post("/", s.clientHandler.CreateClient)
		r.Get("/{id}", s.clientHandler.GetClient)
		r.Put("/{id}", s.clientHandler.UpdateClient)
		r.Delete("/{id}", s.clientHandler.DeleteClient)
	})
}
