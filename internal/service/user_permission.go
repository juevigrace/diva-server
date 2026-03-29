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
	repo *repo.UserPermissionRepository
}

func NewUserPermissionService(repo *repo.UserPermissionRepository) *UserPermissionService {
	return &UserPermissionService{repo: repo}
}

func (s *UserPermissionService) Get(ctx context.Context, userID uuid.UUID) ([]*models.UserPermission, error) {
	return s.repo.GetByUser(ctx, &userID)
}

func (s *UserPermissionService) Create(ctx context.Context, dto *dtos.UserPermissionDto) error {
	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		return err
	}

	var grantedBy *uuid.UUID = nil
	if dto.GrantedBy != "" {
		gb, err := uuid.Parse(dto.GrantedBy)
		if err != nil {
			return err
		}
		grantedBy = &gb
	}

	perm := &models.UserPermission{
		Permission: permissionID,
		UserID:     userID,
		GrantedBy:  grantedBy,
		Granted:    false,
		GrantedAt:  time.Now().UTC().UnixMilli(),
		ExpiresAt:  dto.ExpiresAt,
	}

	return s.repo.Create(ctx, perm)
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

	params := &models.UserPermission{
		Permission: permissionID,
		UserID:     userID,
		Granted:    dto.Granted,
		ExpiresAt:  dto.ExpiresAt,
	}

	return s.repo.Update(ctx, params)
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
	return s.repo.Delete(ctx, &userID, &permissionID)
}

func (s *UserPermissionService) CreateBatch(ctx context.Context, params []*models.UserPermission) error {
	return s.repo.CreateBatch(ctx, params)
}

func (s *UserPermissionService) UpdateBatch(ctx context.Context, params []*models.UserPermission) error {
	return s.repo.UpdateBatch(ctx, params)
}

func (s *UserPermissionService) DeleteBatch(ctx context.Context, userID *uuid.UUID) error {
	return s.repo.DeleteBatch(ctx, userID)
}
