package storage

import (
	"context"

	"github.com/tuanta7/ciam/internal/repository/store"
)

// Storage implements the zitadel op.Storage interface
type Storage struct {
	store.Queries
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Health(ctx context.Context) error {
	return nil
}
