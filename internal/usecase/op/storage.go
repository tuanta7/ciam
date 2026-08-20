package op

import (
	"context"
)

// Storage implements the zitadel op.Storage interface
type Storage struct {
	client ClientUC
}

func NewStorage(
	client ClientUC,
) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) Health(ctx context.Context) error {
	return nil
}
