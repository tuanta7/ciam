package main

import (
	"context"
	"log"
	"os"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/repository"
	"github.com/tuanta7/ciam/internal/transport/rest/handler"
	"github.com/tuanta7/ciam/internal/usecase/client"
	ciamoidc "github.com/tuanta7/ciam/internal/usecase/oidc"
	"github.com/tuanta7/ciam/pkg/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{},
		Action: func(ctx context.Context, command *cli.Command) error {
			cfg := config.LoadConfig()

			executor, err := repository.NewPostgresClient(ctx, cfg.Postgres.DSN)
			if err != nil {
				return err
			}
			defer executor.Close()

			clientRepo := repository.NewClientRepository(executor)
			clientUC := client.NewUseCase(clientRepo)
			clientHandler := handler.NewClientHandler(clientUC)

			provider, err := ciamoidc.NewProvider(cfg.Issuer, cfg.CryptoKey, clientUC)
			if err != nil {
				return err
			}

			server := NewServer(cfg, provider, clientHandler)
			return utils.StartServerWithGracefulShutdown(server)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
