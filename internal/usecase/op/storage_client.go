package op

import (
	"context"

	"github.com/tuanta7/ciam/internal/domain"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type ClientUC interface {
	Get(ctx context.Context, id string) (*domain.Client, error)
}

func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	return s.client.Get(ctx, clientID)
}

func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	client, err := s.client.Get(ctx, clientID)
	if err != nil {
		return err
	}

	if client.Secret != clientSecret {
		return oidc.ErrInvalidClient().WithDescription("invalid client secret")
	}

	return nil
}
