package handler

import (
	"net/http"

	"github.com/tuanta7/ciam/internal/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type OIDCHandler struct {
	provider *oidc.Provider
}

func NewOIDCHandler(provider *oidc.Provider) *OIDCHandler {
	return &OIDCHandler{
		provider: provider,
	}
}

func (h *OIDCHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	op.Authorize(w, r, h.provider)
}

func (h *OIDCHandler) Token(w http.ResponseWriter, r *http.Request) {
	op.TokenExchange(w, r, h.provider)
}

func (h *OIDCHandler) Introspect(w http.ResponseWriter, r *http.Request) {
	op.Introspect(w, r, h.provider)
}
