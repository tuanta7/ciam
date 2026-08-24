package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/repository/orm"
)

type ClientRepository struct {
	exec boil.ContextExecutor
}

func NewClientRepository(exec boil.ContextExecutor) *ClientRepository {
	return &ClientRepository{
		exec: exec,
	}
}

func (r *ClientRepository) Count(ctx context.Context) (int64, error) {
	return orm.Clients().Count(ctx, r.exec)
}

func (r *ClientRepository) List(ctx context.Context, offset, limit int) ([]*domain.Client, error) {
	clients := make([]*domain.Client, 0)

	rows, err := orm.Clients(
		qm.OrderBy(orm.ClientColumns.CreatedAt+" DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	).All(ctx, r.exec)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		clients = append(clients, clientFromRow(row))
	}
	return clients, nil
}

func (r *ClientRepository) Get(ctx context.Context, id string) (*domain.Client, error) {
	row, err := orm.FindClient(ctx, r.exec, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrClientNotFound
	} else if err != nil {
		return nil, err
	}

	return clientFromRow(row), nil
}

func (r *ClientRepository) Create(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	row := rowFromClient(client)
	if err := row.Insert(ctx, r.exec, boil.Infer()); err != nil {
		return nil, err
	}

	// Insert reads the column defaults back, so map the row rather than the
	// argument; the caller may have left defaulted columns empty.
	return clientFromRow(row), nil
}

func (r *ClientRepository) Update(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	row := rowFromClient(client)
	rowsAffected, err := row.Update(ctx, r.exec, boil.Blacklist(
		orm.ClientColumns.Secret,
		orm.ClientColumns.CreatedBy,
		orm.ClientColumns.CreatedAt,
	))
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	// Update only writes the blacklisted-complement columns; reload so the
	// caller sees the untouched columns (secret, created_by, created_at) too.
	if err := row.Reload(ctx, r.exec); err != nil {
		return nil, err
	}

	return clientFromRow(row), nil
}

func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	_, err := (&orm.Client{ID: id}).Delete(ctx, r.exec)
	return err
}

func clientFromRow(row *orm.Client) *domain.Client {
	return &domain.Client{
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

func rowFromClient(c *domain.Client) *orm.Client {
	return &orm.Client{
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
		CreatedAt:                      c.CreatedAt,
		UpdatedAt:                      c.UpdatedAt,
	}
}
