package storage

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type AuthRequestUC interface {
	Create()
}

func (s *Storage) CreateAuthRequest(ctx context.Context, request *oidc.AuthRequest, s2 string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) AuthRequestByID(ctx context.Context, s2 string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) AuthRequestByCode(ctx context.Context, s2 string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SaveAuthCode(ctx context.Context, s3 string, s2 string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteAuthRequest(ctx context.Context, s2 string) error {
	//TODO implement me
	panic("implement me")
}

// AuthenticateAuthRequest marks the auth request as authenticated by subject,
// called by the login handler once credentials have been verified.
func (s *Storage) AuthenticateAuthRequest(ctx context.Context, requestID, subject string) error {
	//TODO implement me
	panic("implement me")
}
