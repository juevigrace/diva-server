package permission

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type PermissionRepo struct {
	queries *db.Queries
}

func NewPermissionRepo(queries *db.Queries) *PermissionRepo {
	return &PermissionRepo{
		queries: queries,
	}
}

func (s *PermissionRepo) List(ctx context.Context, pagination *models.Pagination) ([]*models.Permission, error) {
	rows, err := s.queries.ListPermissions(ctx, db.ListPermissionsParams{
		Limit:  int32(pagination.GetLimit()),
		Offset: int32(pagination.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	permissions := make([]*models.Permission, len(rows))
	for i := range rows {
		permissions[i] = models.PermissionFromDB(&rows[i])
	}

	return permissions, nil
}

func (s *PermissionRepo) Count(ctx context.Context) (int64, error) {
	return s.queries.CountPermissions(ctx)
}

func (s *PermissionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	row, err := s.queries.GetPermissionByID(ctx, models.UUIDPtrToDB(&id))
	if err != nil {
		return nil, err
	}

	return models.PermissionFromDB(&row), nil
}

func (s *PermissionRepo) GetByName(ctx context.Context, action models.PermissionAction) (*models.Permission, error) {
	row, err := s.queries.GetPermissionByName(ctx, action.String())
	if err != nil {
		return nil, err
	}

	return models.PermissionFromDB(&row), nil
}

func (s *PermissionRepo) Create(ctx context.Context, dto *dtos.CreatePermissionDto) error {
	perm := &models.Permission{
		ID:          uuid.New(),
		Name:        dto.Name,
		Description: dto.Description,
		Action:      models.PermissionActionFromString(dto.Action),
		RoleLevel:   models.RoleFromString(dto.RoleLevel),
	}
	return s.queries.CreatePermission(ctx, *perm.DBCreate())
}

func (s *PermissionRepo) Update(ctx context.Context, pid uuid.UUID, dto *dtos.UpdatePermissionDto) error {
	perm := &models.Permission{
		ID:          pid,
		Name:        dto.Name,
		Description: dto.Description,
	}
	return s.queries.UpdatePermission(ctx, *perm.DBUpdate())
}

func (s *PermissionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeletePermission(ctx, models.UUIDPtrToDB(&id))
}

func (s *PermissionRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return s.queries.SoftDeletePermission(ctx, models.UUIDPtrToDB(&id))
}

func (s *PermissionRepo) Restore(ctx context.Context, id uuid.UUID) error {
	return s.queries.RestorePermission(ctx, models.UUIDPtrToDB(&id))
}
