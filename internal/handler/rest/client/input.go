package client

import (
	"fmt"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type ClientInput struct {
	Name                           string   `json:"name" validate:"required"`
	Description                    string   `json:"description"`
	Scopes                         []string `json:"scopes"`
	RedirectURIs                   []string `json:"redirect_uris"`
	PostLogoutRedirectURIs         []string `json:"post_logout_redirect_uris"`
	GrantTypes                     []string `json:"grant_types"`
	ResponseTypes                  []string `json:"response_types"`
	Audiences                      []string `json:"audiences"`
	TokenEndpointAuthMethod        string   `json:"token_endpoint_auth_method" validate:"omitempty,oneof=none client_secret_basic client_secret_post private_key_jwt"`
	ApplicationType                string   `json:"application_type"`
	AccessTokenType                string   `json:"access_token_type"`
	LoginURL                       string   `json:"login_url"`
	IDTokenLifetimeSeconds         int32    `json:"id_token_lifetime_seconds"`
	DevMode                        bool     `json:"dev_mode"`
	ClockSkewSeconds               int32    `json:"clock_skew_seconds"`
	IDTokenUserinfoClaimsAssertion bool     `json:"id_token_userinfo_claims_assertion"`
}

func (in *ClientInput) validate() error {
	if _, err := op.ApplicationTypeString(in.applicationType()); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInvalidClient, err)
	}

	if _, err := op.AccessTokenTypeString(in.accessTokenType()); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInvalidClient, err)
	}

	return nil
}

func (in *ClientInput) scopes() []string {
	if len(in.Scopes) == 0 {
		return []string{oidc.ScopeOpenID}
	}
	return in.Scopes
}

func (in *ClientInput) redirectURIs() []string {
	if in.RedirectURIs == nil {
		return []string{}
	}
	return in.RedirectURIs
}

func (in *ClientInput) postLogoutRedirectURIs() []string {
	if in.PostLogoutRedirectURIs == nil {
		return []string{}
	}
	return in.PostLogoutRedirectURIs
}

func (in *ClientInput) audiences() []string {
	if in.Audiences == nil {
		return []string{}
	}
	return in.Audiences
}

func (in *ClientInput) grantTypes() []string {
	if len(in.GrantTypes) == 0 {
		return []string{string(oidc.GrantTypeCode)}
	}
	return in.GrantTypes
}

func (in *ClientInput) responseTypes() []string {
	if len(in.ResponseTypes) == 0 {
		return []string{string(oidc.ResponseTypeCode)}
	}
	return in.ResponseTypes
}

func (in *ClientInput) tokenEndpointAuthMethod() string {
	if in.TokenEndpointAuthMethod == "" {
		return string(oidc.AuthMethodBasic)
	}
	return in.TokenEndpointAuthMethod
}

func (in *ClientInput) applicationType() string {
	if in.ApplicationType == "" {
		return op.ApplicationTypeWeb.String()
	}
	return in.ApplicationType
}

func (in *ClientInput) accessTokenType() string {
	if in.AccessTokenType == "" {
		return op.AccessTokenTypeBearer.String()
	}
	return in.AccessTokenType
}

func (in *ClientInput) idTokenLifetimeSeconds() int32 {
	if in.IDTokenLifetimeSeconds == 0 {
		return int32(config.DefaultIDTokenLifetime.Seconds())
	}
	return in.IDTokenLifetimeSeconds
}

func (in *ClientInput) toDomain() *domain.Client {
	return &domain.Client{
		Name:                          in.Name,
		Description:                   in.Description,
		Scopes:                        in.scopes(),
		RedirectURIList:               in.redirectURIs(),
		PostLogoutRedirectURIList:     in.postLogoutRedirectURIs(),
		GrantTypeList:                 in.grantTypes(),
		ResponseTypeList:              in.responseTypes(),
		AudienceList:                  in.audiences(),
		TokenEndpointAuthMethod:       in.tokenEndpointAuthMethod(),
		ApplicationTypeName:           in.applicationType(),
		AccessTokenTypeName:           in.accessTokenType(),
		LoginURLTemplate:              in.LoginURL,
		IDTokenLifetimeSeconds:        in.idTokenLifetimeSeconds(),
		DevModeEnabled:                in.DevMode,
		ClockSkewSeconds:              in.ClockSkewSeconds,
		IDTokenUserinfoClaimsAsserted: in.IDTokenUserinfoClaimsAssertion,
	}
}
