package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tuanta7/ciam/internal/repository/store"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var (
	ErrNotFound      = errors.New("client not found")
	ErrInvalidClient = errors.New("invalid client")
)

type Repository interface {
	ListClients(ctx context.Context, arg store.ListClientsParams) ([]store.Client, error)
	GetClient(ctx context.Context, id string) (store.Client, error)
	CreateClient(ctx context.Context, arg store.CreateClientParams) (store.Client, error)
	UpdateClient(ctx context.Context, arg store.UpdateClientParams) (store.Client, error)
	DeleteClient(ctx context.Context, id string) error
}

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

type CreateInput struct {
	ID                             string   `json:"id,omitempty"`
	Name                           string   `json:"name" validate:"required"`
	Description                    string   `json:"description"`
	Secret                         string   `json:"secret"`
	Scopes                         []string `json:"scopes"`
	RedirectURIs                   []string `json:"redirect_uris"`
	PostLogoutRedirectURIs         []string `json:"post_logout_redirect_uris"`
	GrantTypes                     []string `json:"grant_types"`
	ResponseTypes                  []string `json:"response_types"`
	Audiences                      []string `json:"audiences"`
	TokenEndpointAuthMethod        string   `json:"token_endpoint_auth_method"`
	ApplicationType                string   `json:"application_type"`
	AccessTokenType                string   `json:"access_token_type"`
	LoginURL                       string   `json:"login_url"`
	IDTokenLifetimeSeconds         int32    `json:"id_token_lifetime_seconds"`
	DevMode                        bool     `json:"dev_mode"`
	ClockSkewSeconds               int32    `json:"clock_skew_seconds"`
	IDTokenUserinfoClaimsAssertion bool     `json:"id_token_userinfo_claims_assertion"`
	CreatedBy                      string   `json:"created_by"`
	UpdatedBy                      string   `json:"updated_by"`
}

type UpdateInput = CreateInput

func (uc *UseCase) List(ctx context.Context, page, pageSize int32) ([]*Client, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	rows, err := uc.repo.ListClients(ctx, store.ListClientsParams{
		Offset: (page - 1) * pageSize,
		Limit:  pageSize,
	})
	if err != nil {
		return nil, err
	}

	clients := make([]*Client, 0, len(rows))
	for _, row := range rows {
		c, err := fromStore(row)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func (uc *UseCase) Get(ctx context.Context, id string) (*Client, error) {
	row, err := uc.repo.GetClient(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return fromStore(row)
}

func (uc *UseCase) Create(ctx context.Context, in CreateInput) (*Client, error) {
	normalizeInput(&in)
	if in.ID == "" {
		in.ID = uuid.NewString()
	}
	if err := validateEnums(in); err != nil {
		return nil, err
	}

	row, err := uc.repo.CreateClient(ctx, store.CreateClientParams{
		ID:                             in.ID,
		Name:                           in.Name,
		Description:                    in.Description,
		Secret:                         in.Secret,
		Scope:                          mustJSON(in.Scopes),
		RedirectUris:                   mustJSON(in.RedirectURIs),
		PostLogoutRedirectUris:         mustJSON(in.PostLogoutRedirectURIs),
		GrantTypes:                     mustJSON(in.GrantTypes),
		ResponseTypes:                  mustJSON(in.ResponseTypes),
		Audience:                       mustJSON(in.Audiences),
		TokenEndpointAuthMethod:        in.TokenEndpointAuthMethod,
		ApplicationType:                in.ApplicationType,
		AccessTokenType:                in.AccessTokenType,
		LoginUrl:                       in.LoginURL,
		IDTokenLifetimeSeconds:         in.IDTokenLifetimeSeconds,
		DevMode:                        in.DevMode,
		ClockSkewSeconds:               in.ClockSkewSeconds,
		IDTokenUserinfoClaimsAssertion: in.IDTokenUserinfoClaimsAssertion,
		CreatedBy:                      in.CreatedBy,
		UpdatedBy:                      in.UpdatedBy,
	})
	if err != nil {
		return nil, err
	}
	return fromStore(row)
}

func (uc *UseCase) Update(ctx context.Context, id string, in UpdateInput) (*Client, error) {
	normalizeInput(&in)
	if err := validateEnums(in); err != nil {
		return nil, err
	}

	row, err := uc.repo.UpdateClient(ctx, store.UpdateClientParams{
		ID:                             id,
		Name:                           in.Name,
		Description:                    in.Description,
		Secret:                         in.Secret,
		Scope:                          mustJSON(in.Scopes),
		RedirectUris:                   mustJSON(in.RedirectURIs),
		PostLogoutRedirectUris:         mustJSON(in.PostLogoutRedirectURIs),
		GrantTypes:                     mustJSON(in.GrantTypes),
		ResponseTypes:                  mustJSON(in.ResponseTypes),
		Audience:                       mustJSON(in.Audiences),
		TokenEndpointAuthMethod:        in.TokenEndpointAuthMethod,
		ApplicationType:                in.ApplicationType,
		AccessTokenType:                in.AccessTokenType,
		LoginUrl:                       in.LoginURL,
		IDTokenLifetimeSeconds:         in.IDTokenLifetimeSeconds,
		DevMode:                        in.DevMode,
		ClockSkewSeconds:               in.ClockSkewSeconds,
		IDTokenUserinfoClaimsAssertion: in.IDTokenUserinfoClaimsAssertion,
		UpdatedBy:                      in.UpdatedBy,
	})
	if err != nil {
		return nil, mapNotFound(err)
	}
	return fromStore(row)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.DeleteClient(ctx, id)
}

func normalizeInput(in *CreateInput) {
	if len(in.Scopes) == 0 {
		in.Scopes = []string{oidc.ScopeOpenID}
	}
	if len(in.GrantTypes) == 0 {
		in.GrantTypes = []string{string(oidc.GrantTypeCode)}
	}
	if len(in.ResponseTypes) == 0 {
		in.ResponseTypes = []string{string(oidc.ResponseTypeCode)}
	}
	if in.TokenEndpointAuthMethod == "" {
		in.TokenEndpointAuthMethod = string(oidc.AuthMethodBasic)
	}
	if in.ApplicationType == "" {
		in.ApplicationType = op.ApplicationTypeWeb.String()
	}
	if in.AccessTokenType == "" {
		in.AccessTokenType = op.AccessTokenTypeBearer.String()
	}
	if in.IDTokenLifetimeSeconds == 0 {
		in.IDTokenLifetimeSeconds = 3600
	}
	if in.UpdatedBy == "" {
		in.UpdatedBy = in.CreatedBy
	}
}

func validateEnums(in CreateInput) error {
	if _, err := op.ApplicationTypeString(in.ApplicationType); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidClient, err)
	}
	if _, err := op.AccessTokenTypeString(in.AccessTokenType); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidClient, err)
	}
	return nil
}

func fromStore(row store.Client) (*Client, error) {
	return &Client{
		ID:                            row.ID,
		Name:                          row.Name,
		Description:                   row.Description,
		Secret:                        row.Secret,
		Scopes:                        parseJSON(row.Scope),
		RedirectURIList:               parseJSON(row.RedirectUris),
		PostLogoutRedirectURIList:     parseJSON(row.PostLogoutRedirectUris),
		GrantTypeList:                 parseJSON(row.GrantTypes),
		ResponseTypeList:              parseJSON(row.ResponseTypes),
		Audiences:                     parseJSON(row.Audience),
		TokenEndpointAuthMethod:       row.TokenEndpointAuthMethod,
		ApplicationTypeName:           row.ApplicationType,
		AccessTokenTypeName:           row.AccessTokenType,
		LoginURLTemplate:              row.LoginUrl,
		IDTokenLifetimeSeconds:        row.IDTokenLifetimeSeconds,
		DevModeEnabled:                row.DevMode,
		ClockSkewSeconds:              row.ClockSkewSeconds,
		IDTokenUserinfoClaimsAsserted: row.IDTokenUserinfoClaimsAssertion,
		CreatedBy:                     row.CreatedBy,
		UpdatedBy:                     row.UpdatedBy,
		CreatedAt:                     row.CreatedAt.Time,
		UpdatedAt:                     row.UpdatedAt.Time,
	}, nil
}

func parseJSON(raw string) []string {
	if raw == "" {
		return nil
	}
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil
	}
	return values
}

func mustJSON(values []string) string {
	if values == nil {
		values = []string{}
	}
	data, _ := json.Marshal(values)
	return string(data)
}

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
