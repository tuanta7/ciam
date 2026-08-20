package handler

import (
	"net/http"

	ciamop "github.com/tuanta7/ciam/internal/usecase/op"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type OIDCHandler struct {
	provider *ciamop.Provider
}

func NewOIDCHandler(provider *ciamop.Provider) *OIDCHandler {
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
