package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserPermissionService struct {
	repo *repo.UserPermsRepo
}

func NewUserPermissionService(repo *repo.UserPermsRepo) *UserPermissionService {
	return &UserPermissionService{repo: repo}
}

func (s *UserPermissionService) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPermission, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *UserPermissionService) GetOneByUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) (*models.UserPermission, error) {
	return s.repo.GetOneByUser(ctx, userID, permissionID)
}

// TODO: expiration time might better be set here manually instead of being passed from the dto
func (s *UserPermissionService) Create(ctx context.Context, session *models.Session, dto *dtos.UserPermissionDto) error {
	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		return err
	}
	grantedAt := new(int64)
	if dto.Granted {
		*grantedAt = time.Now().UTC().UnixMilli()
	}

	perm := &models.UserPermission{
		Permission: models.Permission{ID: permissionID},
		GrantedBy:  &session.User.ID,
		Granted:    dto.Granted,
		GrantedAt:  grantedAt,
		ExpiresAt:  dto.ExpiresAt,
	}

	return s.repo.Create(ctx, userID, perm)
}

func (s *UserPermissionService) Update(ctx context.Context, dto *dtos.UserPermissionDto) error {
	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		return err
	}
	grantedAt := new(int64)
	if dto.Granted {
		*grantedAt = time.Now().UTC().UnixMilli()
	}

	params := &models.UserPermission{
		Permission: models.Permission{ID: permissionID},
		Granted:    dto.Granted,
		GrantedAt:  grantedAt,
		ExpiresAt:  dto.ExpiresAt,
	}

	return s.repo.Update(ctx, userID, params)
}

func (s *UserPermissionService) Delete(ctx context.Context, dto *dtos.DeleteUserPermissionDto) error {
	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, userID, permissionID)
}

func (s *UserPermissionService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByUser(ctx, userID)
}
