package repo

import (
	"context"
	"time"

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

func (r *UserPermissionRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPermission, error) {
	rows, err := r.queries.GetUserPermissions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	perms := make([]*models.UserPermission, len(rows))
	for i := range rows {
		perms[i] = &models.UserPermission{
			Permission: &models.Permission{ID: rows[i].Permissionid.Bytes},
			GrantedBy:  models.FromUUIDPtr(rows[i].Grantedby),
			Granted:    rows[i].Granted,
			GrantedAt:  models.ToInt64Ptr(rows[i].Grantedat),
			ExpiresAt:  models.ToInt64Ptr(rows[i].Expiresat),
			UpdatedAt:  rows[i].Updatedat.Time.UnixMilli(),
		}
	}
	return perms, nil
}

func (r *UserPermissionRepository) GetByUserAndPermission(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) (*models.UserPermission, error) {
	row, err := r.queries.GetUserPermission(ctx, db.GetUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &models.UserPermission{
		Permission: &models.Permission{ID: row.Permissionid.Bytes},
		GrantedBy:  models.FromUUIDPtr(row.Grantedby),
		Granted:    row.Granted,
		GrantedAt:  models.ToInt64Ptr(row.Grantedat),
		ExpiresAt:  models.ToInt64Ptr(row.Expiresat),
		UpdatedAt:  row.Updatedat.Time.UnixMilli(),
	}, nil
}

func (r *UserPermissionRepository) Grant(ctx context.Context, userID uuid.UUID, up *models.UserPermission) error {
	return r.queries.GrantPermission(ctx, db.GrantPermissionParams{
		PermissionID: pgtype.UUID{Bytes: up.Permission.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		GrantedBy:    models.ToUUIDPtr(up.GrantedBy),
		Granted:      up.Granted,
		GrantedAt:    models.ToTimestamptzPtr(up.GrantedAt),
		ExpiresAt:    models.ToTimestamptzPtr(up.ExpiresAt),
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
}

func (r *UserPermissionRepository) Revoke(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	return r.queries.RevokePermission(ctx, db.RevokePermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserPermissionRepository) UpdateGrant(ctx context.Context, userID uuid.UUID, up *models.UserPermission) error {
	return r.queries.UpdatePermissionGrant(ctx, db.UpdatePermissionGrantParams{
		Granted:      up.Granted,
		GrantedBy:    models.ToUUIDPtr(up.GrantedBy),
		ExpiresAt:    models.ToTimestamptzPtr(up.ExpiresAt),
		PermissionID: pgtype.UUID{Bytes: up.Permission.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserPermissionRepository) Delete(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	return r.queries.DeleteUserPermission(ctx, db.DeleteUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserPermissionRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteAllUserPermissions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
