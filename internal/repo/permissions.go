package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type PermissionsRepo struct {
	queries *db.Queries
}

func NewPermissionsRepo(queries *db.Queries) *PermissionsRepo {
	return &PermissionsRepo{queries: queries}
}

func (r *PermissionsRepo) List(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	rows, err := r.queries.ListPermissions(ctx, db.ListPermissionsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	permissions := make([]*models.Permission, len(rows))
	for i := range rows {
		permissions[i] = &models.Permission{
			ID:          rows[i].ID.Bytes,
			Name:        rows[i].Name,
			Description: rows[i].Description,
			Action:      rows[i].Action,
			RoleLevel:   models.RoleFromDB(rows[i].RoleLevel),
			CreatedAt:   rows[i].CreatedAt.Time.UnixMilli(),
			UpdatedAt:   rows[i].UpdatedAt.Time.UnixMilli(),
			DeletedAt:   models.ToInt64Ptr(rows[i].DeletedAt),
		}
	}
	return permissions, nil
}

func (r *PermissionsRepo) Count(ctx context.Context) (int64, error) {
	return r.queries.CountPermissions(ctx)
}

func (r *PermissionsRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	row, err := r.queries.GetPermissionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.Permission{
		ID:          row.ID.Bytes,
		Name:        row.Name,
		Description: row.Description,
		Action:      row.Action,
		RoleLevel:   models.RoleFromDB(row.RoleLevel),
		CreatedAt:   row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:   row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:   models.ToInt64Ptr(row.DeletedAt),
	}, nil
}

func (r *PermissionsRepo) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	row, err := r.queries.GetPermissionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &models.Permission{
		ID:          row.ID.Bytes,
		Name:        row.Name,
		Description: row.Description,
		Action:      row.Action,
		RoleLevel:   models.RoleFromDB(row.RoleLevel),
		CreatedAt:   row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:   row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:   models.ToInt64Ptr(row.DeletedAt),
	}, nil
}

func (r *PermissionsRepo) Create(ctx context.Context, perm *models.Permission) error {
	return r.queries.CreatePermission(ctx, db.CreatePermissionParams{
		ID:          pgtype.UUID{Bytes: perm.ID, Valid: true},
		Name:        perm.Name,
		Description: perm.Description,
		Action:      perm.Action,
		RoleLevel:   perm.RoleLevel.ToDB(),
	})
}

func (r *PermissionsRepo) Update(ctx context.Context, perm *models.Permission) error {
	return r.queries.UpdatePermission(ctx, db.UpdatePermissionParams{
		Name:        perm.Name,
		Description: perm.Description,
		Action:      perm.Action,
		RoleLevel:   perm.RoleLevel.ToDB(),
		ID:          pgtype.UUID{Bytes: perm.ID, Valid: true},
	})
}

func (r *PermissionsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *PermissionsRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeletePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *PermissionsRepo) Restore(ctx context.Context, id uuid.UUID) error {
	return r.queries.RestorePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
