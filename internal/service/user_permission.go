package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserPermissionService struct {
	queries  *db.Queries
	pService *PermissionService
}

func NewUserPermissionService(queries *db.Queries, pService *PermissionService) *UserPermissionService {
	return &UserPermissionService{
		pService: pService,
		queries:  queries,
	}
}

func (s *UserPermissionService) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPermission, error) {
	rows, err := s.queries.GetUserPermissions(ctx, models.UUIDPtrToDB(&userID))
	if err != nil {
		return nil, err
	}

	perms := make([]models.UserPermission, len(rows))
	for i := range rows {
		perm, err := s.pService.GetByID(ctx, rows[i].PermissionID.Bytes)
		if err != nil {
			continue
		}

		perms[i] = *models.UserPermissionFromDB(&rows[i], perm)
	}

	return perms, nil
}

func (s *UserPermissionService) GetOneByUser(ctx context.Context, userID, permissionID uuid.UUID) (*models.UserPermission, error) {
	row, err := s.queries.GetUserPermission(ctx, db.GetUserPermissionParams{
		PermissionID: models.UUIDPtrToDB(&permissionID),
		UserID:       models.UUIDPtrToDB(&userID),
	})
	if err != nil {
		return nil, err
	}

	dbPerm, err := s.pService.GetByID(ctx, models.DBUUIDToUUID(row.PermissionID))
	if err != nil {
		return nil, err
	}

	return models.UserPermissionFromDB(&row, dbPerm), nil
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

	return s.queries.CreateUserPermission(ctx, *perm.DBCreate(userID))
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

	return s.queries.UpdateUserPermission(ctx, *params.DBUpdate(userID))
}

func (s *UserPermissionService) Delete(ctx context.Context, uid, pid uuid.UUID) error {
	return s.queries.DeleteUserPermission(ctx, db.DeleteUserPermissionParams{
		PermissionID: models.UUIDPtrToDB(&pid),
		UserID:       models.UUIDPtrToDB(&uid),
	})
}

func (s *UserPermissionService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.queries.DeleteAllUserPermissions(ctx, models.UUIDPtrToDB(&userID))
}
