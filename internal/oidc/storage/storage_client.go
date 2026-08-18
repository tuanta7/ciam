package storage

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	return s.client.GetClient(ctx, clientID)
}

func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	client, err := s.client.GetClient(ctx, clientID)
	if err != nil {
		return err
	}

	if client.Secret != clientSecret {
		return oidc.ErrInvalidClient().WithDescription("invalid client secret")
	}

	return nil
}
