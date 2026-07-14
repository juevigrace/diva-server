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

func (s *PermissionRepo) Update(ctx context.Context, pid uuid.UUID, dto *dtos.UpdatePermissionDto) error {
	perm := &models.Permission{
		ID:          pid,
		Name:        dto.Name,
		Description: dto.Description,
	}
	return s.store.UpdatePermission(ctx, perm.DBUpdate())
}

func (s *PermissionRepo) UpdateRoleLevel(ctx context.Context, pid uuid.UUID, role models.Role) error {
	return s.store.UpdatePermissionRoleLevel(ctx, &storage.UpdatePermissionRoleLevelParams{
		ID:        pid,
		RoleLevel: role.ToDB(),
	})
}
