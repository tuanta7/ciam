package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/repository/models"
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
	return models.Clients().Count(ctx, r.exec)
}

func (r *ClientRepository) List(ctx context.Context, offset, limit int) ([]*domain.Client, error) {
	clients := make([]*domain.Client, 0)

	rows, err := models.Clients(
		qm.OrderBy(models.ClientColumns.CreatedAt+" DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	).All(ctx, r.exec)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		clients = append(clients, domain.NewClientFromRow(row))
	}
	return clients, nil
}

func (r *ClientRepository) Get(ctx context.Context, id string) (*domain.Client, error) {
	row, err := models.FindClient(ctx, r.exec, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrClientNotFound
	} else if err != nil {
		return nil, err
	}

	return domain.NewClientFromRow(row), nil
}

func (r *ClientRepository) Create(ctx context.Context, client *models.Client) error {
	return client.Insert(ctx, r.exec, boil.Infer())
}

func (r *ClientRepository) Update(ctx context.Context, client *models.Client) error {
	rowsAffected, err := client.Update(ctx, r.exec, boil.Blacklist(
		models.ClientColumns.Secret,
		models.ClientColumns.CreatedBy,
		models.ClientColumns.CreatedAt,
	))
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	// Update only writes the blacklisted-complement columns; reload so the
	// caller sees the untouched columns (secret, created_by, created_at) too.
	return client.Reload(ctx, r.exec)
}

func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	_, err := (&models.Client{ID: id}).Delete(ctx, r.exec)
	return err
}
