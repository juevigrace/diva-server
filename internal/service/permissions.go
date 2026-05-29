package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
)

type PermissionService struct {
	repo *repo.PermissionsRepo
}

func NewPermissionService(repo *repo.PermissionsRepo) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) List(ctx context.Context, pagination *models.Pagination) ([]*models.Permission, error) {
	return s.repo.List(ctx, pagination.GetLimit(), pagination.GetOffset())
}

func (s *PermissionService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *PermissionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PermissionService) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *PermissionService) Create(ctx context.Context, dto *dtos.CreatePermissionDto) error {
	perm := &models.Permission{
		ID:          uuid.New(),
		Name:        dto.Name,
		Description: dto.Description,
		Action:      models.PermissionActionFromString(dto.Action),
		RoleLevel:   models.RoleFromString(dto.RoleLevel),
	}
	return s.repo.Create(ctx, perm)
}

func (s *PermissionService) Update(ctx context.Context, pid uuid.UUID, dto *dtos.UpdatePermissionDto) error {
	perm := &models.Permission{
		ID:          pid,
		Name:        dto.Name,
		Description: dto.Description,
		Action:      models.PermissionActionFromString(dto.Action),
		RoleLevel:   models.RoleFromString(dto.RoleLevel),
	}
	return s.repo.Update(ctx, perm)
}

func (s *PermissionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *PermissionService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *PermissionService) Restore(ctx context.Context, id uuid.UUID) error {
	return s.repo.Restore(ctx, id)
}
