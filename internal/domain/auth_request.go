package domain

import (
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.AuthRequest = (*AuthRequest)(nil)

type AuthRequest struct {
	ID            string
	ClientID      string
	Subject       string
	RedirectURI   string
	Scopes        []string
	ResponseType  oidc.ResponseType
	ResponseMode  oidc.ResponseMode
	State         string
	Nonce         string
	LoginHint     string
	CodeChallenge *oidc.CodeChallenge
	AuthTime      time.Time
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

// NewAuthRequest maps the parsed /authorize request; the ID and the deadline
// are assigned by the use case.
func NewAuthRequest(request *oidc.AuthRequest) *AuthRequest {
	authRequest := &AuthRequest{
		ClientID:     request.ClientID,
		RedirectURI:  request.RedirectURI,
		Scopes:       request.Scopes,
		ResponseType: request.ResponseType,
		ResponseMode: request.ResponseMode,
		State:        request.State,
		Nonce:        request.Nonce,
		LoginHint:    request.LoginHint,
	}

	if request.CodeChallenge != "" {
		authRequest.CodeChallenge = &oidc.CodeChallenge{
			Challenge: request.CodeChallenge,
			Method:    request.CodeChallengeMethod,
		}
	}

	return authRequest
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
	return a.AuthTime
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
	return a.Subject
}

func (a *AuthRequest) Done() bool {
	return !a.AuthTime.IsZero()
}

// GetAMR reports how the subject authenticated;
// password is the only method so far.
func (a *AuthRequest) GetAMR() []string {
	if a.Done() {
		return []string{"pwd"}
	}
	return nil
}
