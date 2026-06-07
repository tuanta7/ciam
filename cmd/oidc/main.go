package main

import (
	"context"
	"log"
	"os"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/oauth2client"
	"github.com/tuanta7/ciam/internal/repository/postgres"
	"github.com/tuanta7/ciam/internal/repository/store"
	"github.com/tuanta7/ciam/internal/transport/rest"
	"github.com/tuanta7/ciam/internal/transport/rest/handler"
	"github.com/tuanta7/ciam/pkg/graceful"
	"github.com/tuanta7/ciam/pkg/observability"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{},
		Action: func(ctx context.Context, command *cli.Command) error {
			cfg := config.LoadConfig()
			initMonitor(ctx, cfg)

			pool, err := postgres.NewInstrumentedPool(ctx, cfg.Postgres.DSN)
			if err != nil {
				return err
			}
			defer pool.Close()

			pgRepo := store.New(pool)
			clientUC := oauth2client.NewUseCase(pgRepo)
			clientHandler := handler.NewClientHandler(clientUC)

			server := rest.NewServer(cfg.BindAddress, clientHandler)
			return graceful.StartServerWithGracefulShutdown(server)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func initMonitor(ctx context.Context, cfg *config.Config) {
	if !cfg.EnableMetrics {
		observability.InitNoopMeterProvider()
	} else {
		_, err := observability.InitMeterProvider(ctx, cfg.ServiceName, nil)
		if err != nil {
			log.Fatalf("Failed to initialize meter provider: %v", err)
		}
	}

	if !cfg.EnableTracing {
		observability.InitNoopTracerProvider()
	} else {
		_, err := observability.InitTracerProvider(ctx, cfg.ServiceName, nil)
		if err != nil {
			log.Fatalf("Failed to initialize tracer provider: %v", err)
		}
	}
}
