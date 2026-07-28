package postgres

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/tuanta7/ciam/internal/repository/models"
)

type Repository struct {
	exec boil.ContextExecutor
}

func NewClientRepository(exec boil.ContextExecutor) *Repository {
	return &Repository{exec: exec}
}

func (r *Repository) ListClients(ctx context.Context, offset, limit int) (models.ClientSlice, error) {
	return models.Clients(
		qm.OrderBy(models.ClientColumns.CreatedAt+" DESC"),
		qm.Offset(offset),
		qm.Limit(limit),
	).All(ctx, r.exec)
}

func (r *Repository) GetClient(ctx context.Context, id string) (*models.Client, error) {
	return models.FindClient(ctx, r.exec, id)
}

func (r *Repository) CreateClient(ctx context.Context, client *models.Client) error {
	return client.Insert(ctx, r.exec, boil.Infer())
}

func (r *Repository) UpdateClient(ctx context.Context, client *models.Client) error {
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

func (r *Repository) DeleteClient(ctx context.Context, id string) error {
	_, err := (&models.Client{ID: id}).Delete(ctx, r.exec)
	return err
}
