package authrequest

import (
	"context"

	"github.com/google/uuid"
	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/repository"
	"github.com/tuanta7/ciam/pkg/utils"
)

type UseCase struct {
	repo *repository.AuthRequestRepository
}

func NewUseCase(repo *repository.AuthRequestRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, request *domain.AuthRequest) (*domain.AuthRequest, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	request.ID = id.String()
	request.CreatedAt = utils.NowUTC()
	request.ExpiresAt = request.CreatedAt.Add(config.DefaultAuthRequestLifetime)

	if err := uc.repo.Create(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

func (uc *UseCase) Get(ctx context.Context, id string) (*domain.AuthRequest, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *UseCase) Authenticate(ctx context.Context, id, subject string) error {
	return uc.repo.Authenticate(ctx, id, subject, utils.NowUTC())
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
