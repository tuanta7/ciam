package client

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/repository"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo *repository.ClientRepository
}

func NewUseCase(repo *repository.ClientRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) List(ctx context.Context, page, pageSize int32) ([]*domain.Client, error) {

	clients, err := uc.repo.List(ctx, int((page-1)*pageSize), int(pageSize))
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (uc *UseCase) Create(ctx context.Context, client *domain.Client) (*domain.Client, string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, "", err
	}

	var plainSecret string
	if client.TokenEndpointAuthMethod != string(oidc.AuthMethodNone) {
		plainSecret, err = generateSecret()
		if err != nil {
			return nil, "", err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(plainSecret), bcrypt.DefaultCost)
		if err != nil {
			return nil, "", err
		}
		client.Secret = string(hash)
	}

	row := client.ToRow()
	row.ID = id.String()

	if err := uc.repo.Create(ctx, row); err != nil {
		return nil, "", err
	}

	return domain.NewClientFromRow(row), plainSecret, nil
}

func (uc *UseCase) Get(ctx context.Context, id string) (*domain.Client, error) {
	client, err := uc.repo.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrClientNotFound
	} else if err != nil {
		return nil, err
	}
	return client, nil
}

func (uc *UseCase) Update(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	row := client.ToRow()
	err := uc.repo.Update(ctx, row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrClientNotFound
	} else if err != nil {
		return nil, err
	}

	return domain.NewClientFromRow(row), nil
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
