package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/repository"
	clientuc "github.com/tuanta7/ciam/internal/usecase/client"
	"github.com/urfave/cli/v3"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func createClientCommand() *cli.Command {
	return &cli.Command{
		Name:  "create-client",
		Usage: "create a public test client for the authorization code flow",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Value: fmt.Sprintf("test-client-%d", time.Now().Unix())},
			&cli.StringFlag{Name: "redirect-uri", Value: "http://localhost:8080/callback"},
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

			created, secret, err := clientUC.Create(ctx, clientuc.CreateInput{
				Name:                    command.String("name"),
				RedirectURIs:            []string{command.String("redirect-uri")},
				GrantTypes:              []string{string(oidc.GrantTypeCode)},
				ResponseTypes:           []string{string(oidc.ResponseTypeCode)},
				TokenEndpointAuthMethod: string(oidc.AuthMethodNone),
				CreatedBy:               "cli",
			})
			if err != nil {
				return err
			}

			log.Printf("created test client: id=%s redirect_uri=%s\n", created.ID, command.String("redirect-uri"))
			if secret != "" {
				log.Printf("client secret (shown once): %s\n", secret)
			}
			return nil
		},
	}
}
