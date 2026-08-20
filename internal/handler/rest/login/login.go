package login

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/tuanta7/ciam/internal/usecase/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

// loginTemplate is the built-in login UI. A client can point its login_url at
// its own page instead; it only has to POST the same fields back here.
var loginTemplate = template.Must(template.New("login").Parse(``))

var errInvalidCredentials = errors.New("invalid email or password")

type LoginHandler struct {
	provider *oidc.Provider
	issuer   string
}

func (h *LoginHandler) GetLoginForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = loginTemplate.Execute(w, struct {
		AuthRequestID string
		Email         string
		Error         string
	}{r.URL.Query().Get("auth_request_id"), "", ""})
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		h.render(w, r, authRequestID, email, err.Error())
		return
	}

	if err := h.provider.(r.Context(), authRequestID, subject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hand control back to the OP, which turns the request into a code and
	// redirects to the client's redirect_uri.
	ctx := op.ContextWithIssuer(r.Context(), h.issuer)
	callback := op.AuthCallbackURL(h.provider)(ctx, authRequestID)
	http.Redirect(w, r, callback, http.StatusFound)
}

func (h *LoginHandler) render(w http.ResponseWriter, r *http.Request, authRequestID, email, msg string) {

}
