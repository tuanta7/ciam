package storage

import (
	"context"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/op"
)

func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID, clientID string) (*jose.JSONWebKey, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) KeySet(context.Context) ([]op.Key, error) {
	var keys []op.Key

	return keys, nil
}
