package domain

import (
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.AuthRequest = (*AuthRequest)(nil)

type AuthRequest struct {
	subject  string
	authTime time.Time

	ID            string
	ClientID      string
	RedirectURI   string
	Scopes        []string
	ResponseType  oidc.ResponseType
	ResponseMode  oidc.ResponseMode
	State         string
	Nonce         string
	LoginHint     string
	CodeChallenge *oidc.CodeChallenge
	CreatedAt     time.Time
}

func (a *AuthRequest) GetID() string {
	return a.ID
}

func (a *AuthRequest) GetACR() string {
	return ""
}

func (a *AuthRequest) GetAudience() []string {
	return []string{a.ClientID}
}

func (a *AuthRequest) GetAuthTime() time.Time {
	return a.authTime
}

func (a *AuthRequest) GetClientID() string {
	return a.ClientID
}

func (a *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return a.CodeChallenge
}

func (a *AuthRequest) GetNonce() string {
	return a.Nonce
}

func (a *AuthRequest) GetRedirectURI() string {
	return a.RedirectURI
}

func (a *AuthRequest) GetResponseType() oidc.ResponseType {
	return a.ResponseType
}

func (a *AuthRequest) GetResponseMode() oidc.ResponseMode {
	return a.ResponseMode
}

func (a *AuthRequest) GetScopes() []string {
	return a.Scopes
}

func (a *AuthRequest) GetState() string {
	return a.State
}

func (a *AuthRequest) GetSubject() string {
	return a.subject
}

func (a *AuthRequest) Done() bool {
	return !a.authTime.IsZero()
}

// GetAMR reports how the subject authenticated; password is the only method so far.
func (a *AuthRequest) GetAMR() []string {
	if a.Done() {
		return []string{"pwd"}
	}
	return nil
}
