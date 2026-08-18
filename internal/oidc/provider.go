package oidc

import (
	"github.com/tuanta7/ciam/internal/oidc/storage"
	"github.com/tuanta7/ciam/pkg/o11y"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Provider struct {
	*op.Provider
}

func NewProvider(
	clientRepo storage.ClientRepository,
) (*Provider, error) {
	provider, err := op.NewProvider(
		&op.Config{
			CryptoKey: getCryptoKey(),
		},
		storage.NewStorage(clientRepo),
		op.StaticIssuer("ciam"),
		op.WithLogger(o11y.NewLogger("ciam").Logger),
	)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Provider: provider,
	}, nil
}

func getCryptoKey() [32]byte {
	temp := "secret_key_for_crypto_operations"

	var key [32]byte
	copy(key[:], temp)

	return key
}
