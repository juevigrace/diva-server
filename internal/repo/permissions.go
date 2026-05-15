package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type PermissionsRepository struct {
	queries *db.Queries
}

func NewPermissionsRepository(queries *db.Queries) *PermissionsRepository {
	return &PermissionsRepository{queries: queries}
}

func (r *PermissionsRepository) List(ctx context.Context) ([]*models.Permission, error) {
	rows, err := r.queries.ListPermissions(ctx)
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
			RoleLevel:   models.RoleFromDB(rows[i].Rolelevel),
			CreatedAt:   rows[i].Createdat.Time.UnixMilli(),
			UpdatedAt:   rows[i].Updatedat.Time.UnixMilli(),
			DeletedAt:   models.ToInt64Ptr(rows[i].Deletedat),
		}
	}
	return permissions, nil
}

func (r *PermissionsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	row, err := r.queries.GetPermissionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.Permission{
		ID:          row.ID.Bytes,
		Name:        row.Name,
		Description: row.Description,
		Action:      row.Action,
		RoleLevel:   models.RoleFromDB(row.Rolelevel),
		CreatedAt:   row.Createdat.Time.UnixMilli(),
		UpdatedAt:   row.Updatedat.Time.UnixMilli(),
		DeletedAt:   models.ToInt64Ptr(row.Deletedat),
	}, nil
}

func (r *PermissionsRepository) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	row, err := r.queries.GetPermissionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return &models.Permission{
		ID:          row.ID.Bytes,
		Name:        row.Name,
		Description: row.Description,
		Action:      row.Action,
		RoleLevel:   models.RoleFromDB(row.Rolelevel),
		CreatedAt:   row.Createdat.Time.UnixMilli(),
		UpdatedAt:   row.Updatedat.Time.UnixMilli(),
		DeletedAt:   models.ToInt64Ptr(row.Deletedat),
	}, nil
}

func (r *PermissionsRepository) Create(ctx context.Context, perm *models.Permission) error {
	return r.queries.CreatePermission(ctx, db.CreatePermissionParams{
		ID:          pgtype.UUID{Bytes: perm.ID, Valid: true},
		Name:        perm.Name,
		Description: perm.Description,
		Action:      perm.Action,
		RoleLevel:   perm.RoleLevel.ToDB(),
		CreatedAt:   pgtype.Timestamptz{Time: time.UnixMilli(perm.CreatedAt), Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: time.UnixMilli(perm.UpdatedAt), Valid: true},
	})
}

func (r *PermissionsRepository) Update(ctx context.Context, perm *models.Permission) error {
	return r.queries.UpdatePermission(ctx, db.UpdatePermissionParams{
		Name:        perm.Name,
		Description: perm.Description,
		Action:      perm.Action,
		RoleLevel:   perm.RoleLevel.ToDB(),
		ID:          pgtype.UUID{Bytes: perm.ID, Valid: true},
	})
}

func (r *PermissionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *PermissionsRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeletePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *PermissionsRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return r.queries.RestorePermission(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
