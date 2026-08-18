package main

import (
	"context"
	"log"
	"os"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/oauth2client"
	"github.com/tuanta7/ciam/internal/repository/postgres"
	"github.com/tuanta7/ciam/internal/transport/rest"
	"github.com/tuanta7/ciam/internal/transport/rest/handler"
	"github.com/tuanta7/ciam/pkg/graceful"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{},
		Action: func(ctx context.Context, command *cli.Command) error {
			cfg := config.LoadConfig()
			initMonitor(ctx, cfg)

			executor, err := postgres.NewExecutor(ctx, cfg.Postgres.DSN)
			if err != nil {
				return err
			}
			defer executor.Close()

			clientRepo := oauth2client.NewRepository(executor)
			clientUC := oauth2client.NewUseCase(clientRepo)
			clientHandler := handler.NewClientHandler(clientUC)

			server := rest.NewServer(cfg.BindAddress, clientHandler)
			return graceful.StartServerWithGracefulShutdown(server)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
