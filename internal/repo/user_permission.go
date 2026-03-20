package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPermissionRepository struct {
	queries *db.Queries
}

func NewUserPermissionRepository(queries *db.Queries) *UserPermissionRepository {
	return &UserPermissionRepository{queries: queries}
}

func (r *UserPermissionRepository) Create(ctx context.Context, dto *dtos.UserPermissionDto) error {
	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	grantedBy := pgtype.UUID{}
	if dto.GrantedBy != "" {
		gb, err := uuid.Parse(dto.GrantedBy)
		if err != nil {
			return err
		}
		grantedBy = pgtype.UUID{Bytes: gb, Valid: true}
	}
	params := db.CreateUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		GrantedBy:    grantedBy,
		Granted:      dto.Granted,
		GrantedAt:    pgtype.Timestamptz{Time: now, Valid: true},
		ExpiresAt:    models.ToTimestamptzPtr(dto.ExpiresAt),
	}
	return r.queries.CreateUserPermission(ctx, params)
}

func (r *UserPermissionRepository) Update(ctx context.Context, dto *dtos.UserPermissionDto) error {
	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		return err
	}
	params := db.UpdateUserPermissionParams{
		Granted:      dto.Granted,
		ExpiresAt:    models.ToTimestamptzPtr(dto.ExpiresAt),
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	}
	return r.queries.UpdateUserPermission(ctx, params)
}

func (r *UserPermissionRepository) Delete(ctx context.Context, dto *dtos.DeleteUserPermissionDto) error {
	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		return err
	}
	return r.queries.DeleteUserPermission(ctx, db.DeleteUserPermissionParams{
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
}
