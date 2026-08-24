package oidc

import (
	"crypto/sha256"
	"strings"

	"github.com/rs/cors"
	"github.com/tuanta7/ciam/internal/usecase/oidc/storage"
	"github.com/tuanta7/ciam/pkg/otelx"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Provider struct {
	*op.Provider
	Store *storage.Storage
}

func NewProvider(
	issuer, cryptoKey string,
	clientUC storage.ClientUC,
	authRequestUC storage.AuthRequestUC,
) (*Provider, error) {
	opts := []op.Option{
		op.WithCORSOptions(&cors.Options{}),
		op.WithLogger(otelx.NewLogger("ciam").Logger),
	}

	if !strings.HasPrefix(issuer, "https") {
		opts = append(opts, op.WithAllowInsecure())
	}

	storage := storage.New(clientUC, authRequestUC)
	provider, err := op.NewProvider(
		&op.Config{
			// CryptoKey encrypts the authorization code handed to the client.
			CryptoKey: sha256.Sum256([]byte(cryptoKey)),
		},
		storage,
		op.StaticIssuer(issuer),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Provider: provider,
		Store:    storage,
	}, nil
}
