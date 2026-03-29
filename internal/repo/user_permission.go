package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPermissionRepository struct {
	queries *db.Queries
}

func NewUserPermissionRepository(queries *db.Queries) *UserPermissionRepository {
	return &UserPermissionRepository{queries: queries}
}

func (r *UserPermissionRepository) GetByUser(ctx context.Context, id *uuid.UUID) ([]*models.UserPermission, error) {
	rows, err := r.queries.GetUserPermissions(ctx, pgtype.UUID{Bytes: *id, Valid: true})
	if err != nil {
		return nil, err
	}

	perms := make([]*models.UserPermission, len(rows))
	for i, row := range rows {
		var grantedBy *uuid.UUID = nil
		if row.GrantedBy.Valid {
			parsed, err := uuid.ParseBytes(row.GrantedBy.Bytes[:])
			if err != nil {
				return nil, err
			}
			grantedBy = &parsed
		}

		var expiresAt *int64 = nil
		if row.ExpiresAt.Valid {
			var time = row.GrantedAt.Time.UnixMilli()
			expiresAt = &time

		}

		perms[i] = &models.UserPermission{
			Permission: row.PermissionID.Bytes,
			UserID:     row.UserID.Bytes,
			GrantedBy:  grantedBy,
			Granted:    row.Granted,
			GrantedAt:  row.GrantedAt.Time.UnixMilli(),
			ExpiresAt:  expiresAt,
			UpdatedAt:  row.GrantedAt.Time.UnixMilli(),
		}
	}
	return perms, nil
}

func (r *UserPermissionRepository) Create(ctx context.Context, params *models.UserPermission) error {
	return r.queries.CreateUserPermission(ctx, db.CreateUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: params.Permission, Valid: true},
		UserID:       pgtype.UUID{Bytes: params.UserID, Valid: true},
		GrantedBy:    models.ToUUIDPtr(params.GrantedBy),
		Granted:      params.Granted,
		ExpiresAt:    models.ToTimestamptzPtr(params.ExpiresAt),
	})
}

func (r *UserPermissionRepository) Update(ctx context.Context, params *models.UserPermission) error {
	return r.queries.UpdateUserPermission(ctx, db.UpdateUserPermissionParams{
		Granted:      params.Granted,
		ExpiresAt:    models.ToTimestamptzPtr(params.ExpiresAt),
		PermissionID: pgtype.UUID{Bytes: params.Permission, Valid: true},
		UserID:       pgtype.UUID{Bytes: params.UserID, Valid: true},
	})
}

func (r *UserPermissionRepository) Delete(ctx context.Context, userID, permissionID *uuid.UUID) error {
	return r.queries.DeleteUserPermission(ctx, db.DeleteUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: *permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: *userID, Valid: true},
	})
}

func (r *UserPermissionRepository) CreateBatch(ctx context.Context, params []*models.UserPermission) error {
	for _, p := range params {
		if err := r.Create(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserPermissionRepository) UpdateBatch(ctx context.Context, params []*models.UserPermission) error {
	for _, p := range params {
		if err := r.Update(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserPermissionRepository) DeleteBatch(ctx context.Context, userID *uuid.UUID) error {
	return r.queries.DeleteUserPermissions(ctx, pgtype.UUID{Bytes: *userID, Valid: true})
}
