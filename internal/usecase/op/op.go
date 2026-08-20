package op

import (
	"strings"

	"github.com/rs/cors"
	"github.com/tuanta7/ciam/pkg/otelx"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Provider struct {
	*op.Provider
}

func NewProvider(issuer string, clientRepo ClientUC) (*Provider, error) {
	opts := []op.Option{
		op.WithCORSOptions(&cors.Options{}),
		op.WithLogger(otelx.NewLogger("ciam").Logger),
	}

	if !strings.HasPrefix(issuer, "https") {
		opts = append(opts, op.WithAllowInsecure())
	}

	provider, err := op.NewProvider(
		&op.Config{
			CryptoKey: getCryptoKey(),
		},
		NewStorage(clientRepo),
		op.StaticIssuer(issuer),
		opts...,
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
