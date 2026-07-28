package oidc

import (
	"github.com/tuanta7/ciam/internal/oidc/storage"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Provider struct {
	*op.Provider
}

func NewProvider(clientRepo storage.ClientRepository) (*Provider, error) {
	provider, err := op.NewProvider(
		&op.Config{
			CryptoKey: getCryptoKey(),
		},
		storage.NewStorage(clientRepo),
		op.StaticIssuer("ciam"),
	)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Provider: provider,
	}, nil
}

func getCryptoKey() [32]byte {
	return [32]byte{}
}
