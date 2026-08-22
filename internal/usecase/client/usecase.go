package client

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/google/uuid"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/repository"
	"github.com/tuanta7/ciam/internal/repository/models"
	"github.com/tuanta7/ciam/pkg/utils"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo *repository.ClientRepository
}

func NewUseCase(repo *repository.ClientRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) List(ctx context.Context, page, pageSize int32) ([]*domain.Client, int64, error) {
	clients, err := uc.repo.List(ctx, int((page-1)*pageSize), int(pageSize))
	if err != nil {
		return nil, 0, err
	}

	count, err := uc.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return clients, count, nil
}

func (uc *UseCase) Create(ctx context.Context, client *domain.Client) (*domain.Client, string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, "", err
	}

	var plainSecret string
	if client.TokenEndpointAuthMethod != string(oidc.AuthMethodNone) {
		plainSecret, err = generateSecret()
		if err != nil {
			return nil, "", err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(plainSecret), bcrypt.DefaultCost)
		if err != nil {
			return nil, "", err
		}
		client.Secret = string(hash)
	}

	row := rowFromClient(client)
	row.ID = id.String()
	row.CreatedAt = utils.NowUTC()
	row.CreatedBy = ""
	row.UpdatedAt = row.CreatedAt
	row.UpdatedBy = row.CreatedBy

	if err := uc.repo.Create(ctx, row); err != nil {
		return nil, "", err
	}

	return domain.NewClientFromRow(row), plainSecret, nil
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

func (uc *UseCase) Update(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	row := rowFromClient(client)
	row.UpdatedAt = utils.NowUTC()
	row.UpdatedBy = ""

	err := uc.repo.Update(ctx, row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrClientNotFound
	} else if err != nil {
		return nil, err
	}

	return domain.NewClientFromRow(row), nil
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func rowFromClient(c *domain.Client) *models.Client {
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
	}
}
