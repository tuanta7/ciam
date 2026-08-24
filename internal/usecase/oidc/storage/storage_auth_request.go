package storage

import (
	"context"

	"github.com/tuanta7/ciam/internal/domain"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type AuthRequestUC interface {
	Create(ctx context.Context, request *domain.AuthRequest) (*domain.AuthRequest, error)
	Get(ctx context.Context, id string) (*domain.AuthRequest, error)
	Authenticate(ctx context.Context, id, subject string) error
	Delete(ctx context.Context, id string) error
}

// CreateAuthRequest stores the validated /authorize request so the login UI can
// pick it up by ID. userID is only non-empty when an id_token_hint was passed,
// which is not supported yet, so the request always starts unauthenticated.
func (s *Storage) CreateAuthRequest(ctx context.Context, request *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	authRequest, err := s.authRequest.Create(ctx, domain.NewAuthRequest(request))
	if err != nil {
		return nil, err
	}

	return authRequest, nil
}

func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	authRequest, err := s.authRequest.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return authRequest, nil
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteAuthRequest(ctx context.Context, id string) error {
	return s.authRequest.Delete(ctx, id)
}

// AuthenticateAuthRequest marks the auth request as authenticated by subject,
// called by the login handler once credentials have been verified.
func (s *Storage) AuthenticateAuthRequest(ctx context.Context, requestID, subject string) error {
	return s.authRequest.Authenticate(ctx, requestID, subject)
}
