package login

import (
	_ "embed"
	"errors"
	"html/template"
	"net/http"

	"github.com/tuanta7/ciam/internal/usecase/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

//go:embed templates/login.html
var loginTemplateSource string

// loginTemplate is the built-in login UI. A client can point its login_url at
// its own page instead; it only has to POST the same fields back here.
var loginTemplate = template.Must(template.New("login").Parse(loginTemplateSource))

var errInvalidCredentials = errors.New("invalid email or password")

type Handler struct {
	issuer   string
	provider *oidc.Provider
}

func NewHandler(provider *oidc.Provider, issuer string) *Handler {
	return &Handler{
		issuer:   issuer,
		provider: provider,
	}
}

func (h *Handler) GetLoginForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r.URL.Query().Get("auth_request_id"), "", "")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	authRequestID := r.PostForm.Get("auth_request_id")
	email := r.PostForm.Get("email")
	if authRequestID == "" {
		http.Error(w, "missing auth_request_id", http.StatusBadRequest)
		return
	}

	subject, err := h.authenticate(email, r.PostForm.Get("password"))
	if err != nil {
		h.render(w, authRequestID, email, err.Error())
		return
	}

	ctx := r.Context()
	if err := h.provider.Store.AuthenticateAuthRequest(ctx, authRequestID, subject); err != nil {
		h.render(w, authRequestID, email, err.Error())
		return
	}

	// Hand control back to the OP, which turns the request into a code and
	// redirects to the client's redirect_uri.
	ctx = op.ContextWithIssuer(ctx, h.issuer)
	callback := op.AuthCallbackURL(h.provider)(ctx, authRequestID)
	http.Redirect(w, r, callback, http.StatusFound)
}

func (h *Handler) render(w http.ResponseWriter, authRequestID, email, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = loginTemplate.Execute(w, struct {
		AuthRequestID string
		Email         string
		Error         string
	}{authRequestID, email, errMsg})
}

func (h *Handler) authenticate(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errInvalidCredentials
	}

	return "user", nil
}
