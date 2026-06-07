package storage

import (
	"context"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	//TODO implement me
	panic("implement me")
}

// SetUserinfoFromScopes is deprecated and should have an empty implementation for now.
func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID, clientID string, scopes []string) error {
	return nil
}
