package storage

import (
	"context"

	"github.com/tuanta7/ciam/internal/repository/models"
)

type ClientRepository interface {
	GetClient(ctx context.Context, id string) (*models.Client, error)
}

// Storage implements the zitadel op.Storage interface
type Storage struct {
	clientRepo ClientRepository
}

func NewStorage(clientRepo ClientRepository) *Storage {
	return &Storage{clientRepo: clientRepo}
}

func (s *Storage) GetClient(ctx context.Context, id string) (*models.Client, error) {
	return s.clientRepo.GetClient(ctx, id)
}

func (s *Storage) Health(ctx context.Context) error {
	return nil
}
