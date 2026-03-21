package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/storage/db"
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
	grantedBy := pgtype.UUID{}
	if dto.GrantedBy != "" {
		gb, err := uuid.Parse(dto.GrantedBy)
		if err != nil {
			return err
		}
		grantedBy = pgtype.UUID{Bytes: gb, Valid: true}
	}

	params := &db.CreateUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		GrantedBy:    grantedBy,
		Granted:      dto.Granted,
		GrantedAt:    pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		ExpiresAt:    models.ToTimestamptzPtr(dto.ExpiresAt),
	}

	return s.repo.Create(ctx, params)
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
	params := &db.UpdateUserPermissionParams{
		Granted:      dto.Granted,
		ExpiresAt:    models.ToTimestamptzPtr(dto.ExpiresAt),
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
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
	return s.repo.Delete(ctx, userID, permissionID)
}
