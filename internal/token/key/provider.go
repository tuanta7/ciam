package key

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/op"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (m *Provider) KeySet(ctx context.Context) ([]op.Key, error) {
	return []op.Key{}, nil
}
