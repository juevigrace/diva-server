package permission

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage"
)

type PermissionRepo struct {
	store storage.PermissionStore
}

func NewPermissionRepo(store storage.PermissionStore) *PermissionRepo {
	return &PermissionRepo{
		store: store,
	}
}

func (s *PermissionRepo) List(ctx context.Context, pagination *models.Pagination) ([]*models.Permission, error) {
	rows, err := s.store.ListPermissions(ctx, &storage.ListPermissionsParams{
		Limit:  int64(pagination.GetLimit()),
		Offset: int64(pagination.GetOffset()),
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
	return s.store.CountPermissions(ctx)
}

func (s *PermissionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	row, err := s.store.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return models.PermissionFromDB(row), nil
}

func (s *PermissionRepo) GetByName(ctx context.Context, action models.PermissionAction) (*models.Permission, error) {
	row, err := s.store.GetPermissionByName(ctx, action.String())
	if err != nil {
		return nil, err
	}

	return models.PermissionFromDB(row), nil
}

func (s *PermissionRepo) Create(ctx context.Context, dto *dtos.CreatePermissionDto) error {
	perm := &models.Permission{
		ID:          uuid.New(),
		Name:        dto.Name,
		Description: dto.Description,
		Action:      models.PermissionActionFromString(dto.Action),
		RoleLevel:   models.RoleFromString(dto.RoleLevel),
	}
	return s.store.CreatePermission(ctx, perm.DBCreate())
}

func (s *PermissionRepo) Update(ctx context.Context, pid uuid.UUID, dto *dtos.UpdatePermissionDto) error {
	perm := &models.Permission{
		ID:          pid,
		Name:        dto.Name,
		Description: dto.Description,
	}
	return s.store.UpdatePermission(ctx, perm.DBUpdate())
}

func (s *PermissionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return s.store.DeletePermission(ctx, id)
}

func (s *PermissionRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return s.store.SoftDeletePermission(ctx, id)
}

func (s *PermissionRepo) Restore(ctx context.Context, id uuid.UUID) error {
	return s.store.RestorePermission(ctx, id)
}
