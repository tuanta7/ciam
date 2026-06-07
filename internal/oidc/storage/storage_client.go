package storage

import (
	"context"

	"github.com/tuanta7/ciam/internal/oauth2client"
	"github.com/zitadel/oidc/v3/pkg/op"
)

func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return oauth2client.NewClientFromStore(client), nil
}
