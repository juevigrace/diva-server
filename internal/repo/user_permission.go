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

		var grandedAt *int64 = nil
		if row.GrantedAt.Valid {
			var time = row.GrantedAt.Time.UnixMilli()
			grandedAt = &time
		}

		var expiresAt *int64 = nil
		if row.ExpiresAt.Valid {
			var time = row.GrantedAt.Time.UnixMilli()
			expiresAt = &time

		}

		perms[i] = &models.UserPermission{
			Permission: row.PermissionID.Bytes,
			User:       row.UserID.Bytes,
			GrantedBy:  grantedBy,
			Granted:    row.Granted,
			GrantedAt:  grandedAt,
			ExpiresAt:  expiresAt,
			UpdatedAt:  row.GrantedAt.Time.UnixMilli(),
		}
	}
	return perms, nil
}

func (r *UserPermissionRepository) Create(ctx context.Context, params *db.CreateUserPermissionParams) error {
	return r.queries.CreateUserPermission(ctx, *params)
}

func (r *UserPermissionRepository) Update(ctx context.Context, params *db.UpdateUserPermissionParams) error {
	return r.queries.UpdateUserPermission(ctx, *params)
}

func (r *UserPermissionRepository) Delete(ctx context.Context, userID, permissionID uuid.UUID) error {
	return r.queries.DeleteUserPermission(ctx, db.DeleteUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
}
