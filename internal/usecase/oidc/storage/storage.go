package storage

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.Storage = &Storage{}

// Storage implements the zitadel op.Storage interface
type Storage struct {
	client      ClientUC
	authRequest AuthRequestUC
}

func New(c ClientUC, ar AuthRequestUC) *Storage {
	return &Storage{
		client:      c,
		authRequest: ar,
	}
}

func (s *Storage) Health(ctx context.Context) error {
	return nil
}
