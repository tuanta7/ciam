package main

import (
	"context"
	"log"
	"os"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/repository"
	"github.com/tuanta7/ciam/internal/transport/rest"
	"github.com/tuanta7/ciam/internal/transport/rest/handler"
	"github.com/tuanta7/ciam/internal/usecase/client"
	"github.com/tuanta7/ciam/pkg/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{},
		Action: func(ctx context.Context, command *cli.Command) error {
			cfg := config.LoadConfig()
			initMonitor(ctx, cfg)

			executor, err := repository.NewPostgresClient(ctx, cfg.Postgres.DSN)
			if err != nil {
				return err
			}
			defer executor.Close()

			clientRepo := repository.NewClientRepository(executor)
			clientUC := client.NewUseCase(clientRepo)
			clientHandler := handler.NewClientHandler(clientUC)

			server := rest.NewServer(cfg.BindAddress, clientHandler)
			return utils.StartServerWithGracefulShutdown(server)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
