package storage

import (
	"context"
)

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
