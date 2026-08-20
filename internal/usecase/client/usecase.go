package client

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/google/uuid"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/repository"
	"github.com/tuanta7/ciam/internal/repository/models"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type UseCase struct {
	repo *repository.ClientRepository
}

func NewUseCase(repo *repository.ClientRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) List(ctx context.Context, page, pageSize int32) ([]*domain.Client, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	clients, err := uc.repo.List(ctx, int((page-1)*pageSize), int(pageSize))
	if err != nil {
		return nil, err
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

func (uc *UseCase) Create(ctx context.Context, in CreateInput) (*domain.Client, error) {
	in.normalize()
	if err := in.validate(); err != nil {
		return nil, err
	}

	client := &models.Client{
		ID:                             in.ID,
		Name:                           in.Name,
		Description:                    in.Description,
		Secret:                         "",
		Scope:                          types.StringArray(in.Scopes),
		RedirectUris:                   types.StringArray(in.RedirectURIs),
		PostLogoutRedirectUris:         types.StringArray(in.PostLogoutRedirectURIs),
		GrantTypes:                     types.StringArray(in.GrantTypes),
		ResponseTypes:                  types.StringArray(in.ResponseTypes),
		Audience:                       types.StringArray(in.Audiences),
		TokenEndpointAuthMethod:        in.TokenEndpointAuthMethod,
		ApplicationType:                in.ApplicationType,
		AccessTokenType:                in.AccessTokenType,
		LoginURL:                       in.LoginURL,
		IDTokenLifetimeSeconds:         int(in.IDTokenLifetimeSeconds),
		DevMode:                        in.DevMode,
		ClockSkewSeconds:               int(in.ClockSkewSeconds),
		IDTokenUserinfoClaimsAssertion: in.IDTokenUserinfoClaimsAssertion,
		CreatedBy:                      in.CreatedBy,
		UpdatedBy:                      in.UpdatedBy,
	}
	if err := uc.repo.Create(ctx, client); err != nil {
		return nil, err
	}

	return domain.NewClientFromRow(client), nil
}

func (uc *UseCase) Get(ctx context.Context, id string) (*domain.Client, error) {
	client, err := uc.repo.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrClientNotFound
	} else if err != nil {
		return nil, err
	}
	return client, nil
}

type UpdateInput = CreateInput

func (uc *UseCase) Update(ctx context.Context, id string, in UpdateInput) (*domain.Client, error) {
	in.normalize()
	if err := in.validate(); err != nil {
		return nil, err
	}

	client := &models.Client{
		ID:                             id,
		Name:                           in.Name,
		Description:                    in.Description,
		Scope:                          types.StringArray(in.Scopes),
		RedirectUris:                   types.StringArray(in.RedirectURIs),
		PostLogoutRedirectUris:         types.StringArray(in.PostLogoutRedirectURIs),
		GrantTypes:                     types.StringArray(in.GrantTypes),
		ResponseTypes:                  types.StringArray(in.ResponseTypes),
		Audience:                       types.StringArray(in.Audiences),
		TokenEndpointAuthMethod:        in.TokenEndpointAuthMethod,
		ApplicationType:                in.ApplicationType,
		AccessTokenType:                in.AccessTokenType,
		LoginURL:                       in.LoginURL,
		IDTokenLifetimeSeconds:         int(in.IDTokenLifetimeSeconds),
		DevMode:                        in.DevMode,
		ClockSkewSeconds:               int(in.ClockSkewSeconds),
		IDTokenUserinfoClaimsAssertion: in.IDTokenUserinfoClaimsAssertion,
		UpdatedBy:                      in.UpdatedBy,
	}

	err := uc.repo.Update(ctx, client)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrClientNotFound
	} else if err != nil {
		return nil, err
	}

	return domain.NewClientFromRow(client), nil
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (in *CreateInput) normalize() {
	if in.ID == "" {
		in.ID = uuid.NewString()
	}

	if in.RedirectURIs == nil {
		in.RedirectURIs = []string{}
	}

	if in.PostLogoutRedirectURIs == nil {
		in.PostLogoutRedirectURIs = []string{}
	}

	if in.Audiences == nil {
		in.Audiences = []string{}
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
		return fmt.Errorf("%w: %v", domain.ErrInvalidClient, err)
	}

	if _, err := op.AccessTokenTypeString(in.AccessTokenType); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInvalidClient, err)
	}

	return nil
}
