package domain

import (
	"fmt"
	"slices"
	"time"

	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/internal/repository/models"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.Client = (*Client)(nil)

type Client struct {
	ID                            string    `json:"id"`
	Name                          string    `json:"name"`
	Description                   string    `json:"description"`
	Secret                        string    `json:"-"`
	Scopes                        []string  `json:"scopes"`
	RedirectURIList               []string  `json:"redirect_uris"`
	PostLogoutRedirectURIList     []string  `json:"post_logout_redirect_uris"`
	GrantTypeList                 []string  `json:"grant_types"`
	ResponseTypeList              []string  `json:"response_types"`
	AudienceList                  []string  `json:"audiences"`
	TokenEndpointAuthMethod       string    `json:"token_endpoint_auth_method"`
	ApplicationTypeName           string    `json:"application_type"`
	AccessTokenTypeName           string    `json:"access_token_type"`
	LoginURLTemplate              string    `json:"login_url"`
	IDTokenLifetimeSeconds        int32     `json:"id_token_lifetime_seconds"`
	DevModeEnabled                bool      `json:"dev_mode"`
	ClockSkewSeconds              int32     `json:"clock_skew_seconds"`
	IDTokenUserinfoClaimsAsserted bool      `json:"id_token_userinfo_claims_assertion"`
	CreatedBy                     string    `json:"created_by"`
	UpdatedBy                     string    `json:"updated_by"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) RedirectURIs() []string {
	return c.RedirectURIList
}

func (c *Client) PostLogoutRedirectURIs() []string {
	return c.PostLogoutRedirectURIList
}

func (c *Client) ApplicationType() op.ApplicationType {
	appType, err := op.ApplicationTypeString(c.ApplicationTypeName)
	if err != nil {
		return op.ApplicationTypeWeb
	}
	return appType
}

func (c *Client) AuthMethod() oidc.AuthMethod {
	return oidc.AuthMethod(c.TokenEndpointAuthMethod)
}

func (c *Client) ResponseTypes() []oidc.ResponseType {
	values := make([]oidc.ResponseType, 0, len(c.ResponseTypeList))
	for _, responseType := range c.ResponseTypeList {
		values = append(values, oidc.ResponseType(responseType))
	}
	return values
}

func (c *Client) GrantTypes() []oidc.GrantType {
	values := make([]oidc.GrantType, 0, len(c.GrantTypeList))
	for _, grantType := range c.GrantTypeList {
		values = append(values, oidc.GrantType(grantType))
	}
	return values
}

func (c *Client) LoginURL(authReqID string) string {
	if c.LoginURLTemplate == "" {
		return fmt.Sprintf(config.DefaultLoginURLTemplate, authReqID)
	}
	return fmt.Sprintf(c.LoginURLTemplate, authReqID)
}

func (c *Client) AccessTokenType() op.AccessTokenType {
	tokenType, err := op.AccessTokenTypeString(c.AccessTokenTypeName)
	if err != nil {
		// default to opaque token
		return op.AccessTokenTypeBearer
	}

	return tokenType
}

func (c *Client) IDTokenLifetime() time.Duration {
	return time.Duration(c.IDTokenLifetimeSeconds) * time.Second
}

func (c *Client) DevMode() bool {
	return c.DevModeEnabled
}

func (c *Client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return c.restrictScopes
}

func (c *Client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return c.restrictScopes
}

func (c *Client) IsScopeAllowed(scope string) bool {
	return slices.Contains(c.Scopes, scope)
}

func (c *Client) IDTokenUserinfoClaimsAssertion() bool {
	return c.IDTokenUserinfoClaimsAsserted
}

func (c *Client) ClockSkew() time.Duration {
	return time.Duration(c.ClockSkewSeconds) * time.Second
}

func (c *Client) restrictScopes(scopes []string) []string {
	allowed := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		if c.IsScopeAllowed(scope) {
			allowed = append(allowed, scope)
		}
	}
	return allowed
}

func (c *Client) ToRow() *models.Client {
	return &models.Client{
		ID:                             c.ID,
		Name:                           c.Name,
		Description:                    c.Description,
		Secret:                         c.Secret,
		Scope:                          types.StringArray(c.Scopes),
		RedirectUris:                   types.StringArray(c.RedirectURIList),
		PostLogoutRedirectUris:         types.StringArray(c.PostLogoutRedirectURIList),
		GrantTypes:                     types.StringArray(c.GrantTypeList),
		ResponseTypes:                  types.StringArray(c.ResponseTypeList),
		Audience:                       types.StringArray(c.AudienceList),
		TokenEndpointAuthMethod:        c.TokenEndpointAuthMethod,
		ApplicationType:                c.ApplicationTypeName,
		AccessTokenType:                c.AccessTokenTypeName,
		LoginURL:                       c.LoginURLTemplate,
		IDTokenLifetimeSeconds:         int(c.IDTokenLifetimeSeconds),
		DevMode:                        c.DevModeEnabled,
		ClockSkewSeconds:               int(c.ClockSkewSeconds),
		IDTokenUserinfoClaimsAssertion: c.IDTokenUserinfoClaimsAsserted,
		CreatedBy:                      c.CreatedBy,
		UpdatedBy:                      c.UpdatedBy,
	}
}

func NewClientFromRow(row *models.Client) *Client {
	return &Client{
		ID:                            row.ID,
		Name:                          row.Name,
		Description:                   row.Description,
		Secret:                        row.Secret,
		Scopes:                        row.Scope,
		RedirectURIList:               row.RedirectUris,
		PostLogoutRedirectURIList:     row.PostLogoutRedirectUris,
		GrantTypeList:                 row.GrantTypes,
		ResponseTypeList:              row.ResponseTypes,
		AudienceList:                  row.Audience,
		TokenEndpointAuthMethod:       row.TokenEndpointAuthMethod,
		ApplicationTypeName:           row.ApplicationType,
		AccessTokenTypeName:           row.AccessTokenType,
		LoginURLTemplate:              row.LoginURL,
		IDTokenLifetimeSeconds:        int32(row.IDTokenLifetimeSeconds),
		DevModeEnabled:                row.DevMode,
		ClockSkewSeconds:              int32(row.ClockSkewSeconds),
		IDTokenUserinfoClaimsAsserted: row.IDTokenUserinfoClaimsAssertion,
		CreatedBy:                     row.CreatedBy,
		UpdatedBy:                     row.UpdatedBy,
		CreatedAt:                     row.CreatedAt,
		UpdatedAt:                     row.UpdatedAt,
	}
}
