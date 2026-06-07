package client

import (
	"context"
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
		clients = append(clients, NewClientFromStore(row))
	}
	return clients, nil
}

type CreateInput struct {
	ID                             string   `json:"id,omitempty"`
	Name                           string   `json:"name" validate:"required"`
	Description                    string   `json:"description"`
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

func (uc *UseCase) Create(ctx context.Context, in CreateInput) (*Client, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}

	in.normalize()

	row, err := uc.repo.CreateClient(ctx, store.CreateClientParams{
		ID:                             in.ID,
		Name:                           in.Name,
		Description:                    in.Description,
		Secret:                         "",
		Scope:                          in.Scopes,
		RedirectUris:                   in.RedirectURIs,
		PostLogoutRedirectUris:         in.PostLogoutRedirectURIs,
		GrantTypes:                     in.GrantTypes,
		ResponseTypes:                  in.ResponseTypes,
		Audience:                       in.Audiences,
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

	return NewClientFromStore(row), nil
}

func (uc *UseCase) Get(ctx context.Context, id string) (*Client, error) {
	row, err := uc.repo.GetClient(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return NewClientFromStore(row), nil
}

type UpdateInput = CreateInput

func (uc *UseCase) Update(ctx context.Context, id string, in UpdateInput) (*Client, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}

	in.normalize()
	row, err := uc.repo.UpdateClient(ctx, store.UpdateClientParams{
		ID:                             id,
		Name:                           in.Name,
		Description:                    in.Description,
		Scope:                          in.Scopes,
		RedirectUris:                   in.RedirectURIs,
		PostLogoutRedirectUris:         in.PostLogoutRedirectURIs,
		GrantTypes:                     in.GrantTypes,
		ResponseTypes:                  in.ResponseTypes,
		Audience:                       in.Audiences,
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

	return NewClientFromStore(row), nil
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.DeleteClient(ctx, id)
}

func (in *CreateInput) normalize() {
	if in.ID == "" {
		in.ID = uuid.NewString()
	}

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

func (in *CreateInput) validate() error {
	if _, err := op.ApplicationTypeString(in.ApplicationType); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidClient, err)
	}

	if _, err := op.AccessTokenTypeString(in.AccessTokenType); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidClient, err)
	}

	return nil
}

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
