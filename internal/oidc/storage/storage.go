package storage

import (
	"context"

	"github.com/tuanta7/ciam/internal/oauth2client"
)

type ClientRepository interface {
	GetClient(ctx context.Context, id string) (*oauth2client.Client, error)
}

// Storage implements the zitadel op.Storage interface
type Storage struct {
	client ClientRepository
}

func NewStorage(
	client ClientRepository,
) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) Health(ctx context.Context) error {
	return nil
}
