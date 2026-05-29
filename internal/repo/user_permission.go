package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPermsRepo struct {
	queries *db.Queries
}

func NewUserPermsRepo(queries *db.Queries) *UserPermsRepo {
	return &UserPermsRepo{queries: queries}
}

func (r *UserPermsRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPermission, error) {
	rows, err := r.queries.GetUserPermissions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	perms := make([]models.UserPermission, len(rows))
	for i := range rows {
		perms[i] = models.UserPermission{
			Permission: models.Permission{ID: rows[i].PermissionID.Bytes},
			GrantedBy:  models.FromUUIDPtr(rows[i].GrantedBy),
			Granted:    rows[i].Granted,
			GrantedAt:  models.ToInt64Ptr(rows[i].GrantedAt),
			ExpiresAt:  models.ToInt64Ptr(rows[i].ExpiresAt),
			UpdatedAt:  rows[i].UpdatedAt.Time.UnixMilli(),
		}
	}
	return perms, nil
}

func (r *UserPermsRepo) GetOneByUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) (*models.UserPermission, error) {
	row, err := r.queries.GetUserPermission(ctx, db.GetUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &models.UserPermission{
		Permission: models.Permission{ID: row.PermissionID.Bytes},
		GrantedBy:  models.FromUUIDPtr(row.GrantedBy),
		Granted:    row.Granted,
		GrantedAt:  models.ToInt64Ptr(row.GrantedAt),
		ExpiresAt:  models.ToInt64Ptr(row.ExpiresAt),
		UpdatedAt:  row.UpdatedAt.Time.UnixMilli(),
	}, nil
}

func (r *UserPermsRepo) Create(ctx context.Context, userID uuid.UUID, up *models.UserPermission) error {
	return r.queries.CreateUserPermission(ctx, db.CreateUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: up.Permission.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		GrantedBy:    models.ToUUIDPtr(up.GrantedBy),
		Granted:      up.Granted,
		GrantedAt:    models.ToTimestamptzPtr(up.GrantedAt),
		ExpiresAt:    models.ToTimestamptzPtr(up.ExpiresAt),
	})
}

func (r *UserPermsRepo) Update(ctx context.Context, userID uuid.UUID, up *models.UserPermission) error {
	return r.queries.UpdateUserPermission(ctx, db.UpdateUserPermissionParams{
		Granted:      up.Granted,
		GrantedAt:    models.ToTimestamptzPtr(up.GrantedAt),
		ExpiresAt:    models.ToTimestamptzPtr(up.ExpiresAt),
		PermissionID: pgtype.UUID{Bytes: up.Permission.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserPermsRepo) Delete(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	return r.queries.DeleteUserPermission(ctx, db.DeleteUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserPermsRepo) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteAllUserPermissions(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}
