package main

import (
	"context"
	"log"
	"os"

	"github.com/tuanta7/ciam/internal/config"
	clientrest "github.com/tuanta7/ciam/internal/handler/rest/client"
	"github.com/tuanta7/ciam/internal/handler/rest/login"
	"github.com/tuanta7/ciam/internal/repository"
	authrequestuc "github.com/tuanta7/ciam/internal/usecase/authrequest"
	clientuc "github.com/tuanta7/ciam/internal/usecase/client"
	"github.com/tuanta7/ciam/internal/usecase/oidc"
	"github.com/tuanta7/ciam/pkg/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			createClientCommand(),
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			cfg := config.LoadConfig()

			executor, err := repository.NewPostgresClient(ctx, cfg.Postgres.DSN)
			if err != nil {
				return err
			}
			defer executor.Close()

			clientRepo := repository.NewClientRepository(executor)
			clientUC := clientuc.NewUseCase(clientRepo)

			authRequestRepo := repository.NewAuthRequestRepository(executor)
			authRequestUC := authrequestuc.NewUseCase(authRequestRepo)

			clientHandler := clientrest.NewHandler(clientUC)

			provider, err := oidc.NewProvider(
				cfg.Issuer,
				cfg.CryptoKey,
				clientUC,
				authRequestUC,
			)
			if err != nil {
				return err
			}

			authenticator := login.NewHandler(provider, cfg.Issuer)
			server := NewServer(cfg, provider, authenticator, clientHandler)
			return utils.StartServerWithGracefulShutdown(server)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
