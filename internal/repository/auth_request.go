package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/repository/orm"
	"github.com/tuanta7/ciam/pkg/utils"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type AuthRequestRepository struct {
	exec boil.ContextExecutor
}

func NewAuthRequestRepository(exec boil.ContextExecutor) *AuthRequestRepository {
	return &AuthRequestRepository{
		exec: exec,
	}
}

func (r *AuthRequestRepository) Create(ctx context.Context, request *domain.AuthRequest) error {
	return rowFromAuthRequest(request).Insert(ctx, r.exec, boil.Infer())
}

func (r *AuthRequestRepository) Get(ctx context.Context, id string) (*domain.AuthRequest, error) {
	row, err := orm.AuthRequests(
		orm.AuthRequestWhere.ID.EQ(id),
		orm.AuthRequestWhere.ExpiresAt.GT(utils.NowUTC()),
	).One(ctx, r.exec)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAuthRequestNotFound
	} else if err != nil {
		return nil, err
	}

	return authRequestFromRow(row), nil
}

func (r *AuthRequestRepository) Authenticate(ctx context.Context, id, subject string, authTime time.Time) error {
	rowsAffected, err := orm.AuthRequests(
		orm.AuthRequestWhere.ID.EQ(id),
		orm.AuthRequestWhere.ExpiresAt.GT(utils.NowUTC()),
	).UpdateAll(ctx, r.exec, orm.M{
		orm.AuthRequestColumns.Subject:  subject,
		orm.AuthRequestColumns.AuthTime: null.TimeFrom(authTime),
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrAuthRequestNotFound
	}

	return nil
}

func (r *AuthRequestRepository) Delete(ctx context.Context, id string) error {
	_, err := (&orm.AuthRequest{ID: id}).Delete(ctx, r.exec)
	return err
}

func rowFromAuthRequest(a *domain.AuthRequest) *orm.AuthRequest {
	row := &orm.AuthRequest{
		ID:           a.ID,
		ClientID:     a.ClientID,
		Subject:      a.Subject,
		RedirectURI:  a.RedirectURI,
		Scopes:       types.StringArray(a.Scopes),
		ResponseType: string(a.ResponseType),
		ResponseMode: string(a.ResponseMode),
		State:        a.State,
		Nonce:        a.Nonce,
		LoginHint:    a.LoginHint,
		AuthTime:     null.NewTime(a.AuthTime, !a.AuthTime.IsZero()),
		CreatedAt:    a.CreatedAt,
		ExpiresAt:    a.ExpiresAt,
	}

	if a.CodeChallenge != nil {
		row.CodeChallenge = a.CodeChallenge.Challenge
		row.CodeChallengeMethod = string(a.CodeChallenge.Method)
	}

	return row
}

func authRequestFromRow(row *orm.AuthRequest) *domain.AuthRequest {
	authRequest := &domain.AuthRequest{
		ID:           row.ID,
		ClientID:     row.ClientID,
		Subject:      row.Subject,
		RedirectURI:  row.RedirectURI,
		Scopes:       row.Scopes,
		ResponseType: oidc.ResponseType(row.ResponseType),
		ResponseMode: oidc.ResponseMode(row.ResponseMode),
		State:        row.State,
		Nonce:        row.Nonce,
		LoginHint:    row.LoginHint,
		AuthTime:     row.AuthTime.Time,
		CreatedAt:    row.CreatedAt,
		ExpiresAt:    row.ExpiresAt,
	}

	if row.CodeChallenge != "" {
		authRequest.CodeChallenge = &oidc.CodeChallenge{
			Challenge: row.CodeChallenge,
			Method:    oidc.CodeChallengeMethod(row.CodeChallengeMethod),
		}
	}

	return authRequest
}
