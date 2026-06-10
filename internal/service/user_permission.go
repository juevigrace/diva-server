package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
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

func (s *UserPermissionService) GetOneByPermID(ctx context.Context, userID, permissionID uuid.UUID) (*models.UserPermission, error) {
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

func (s *UserPermissionService) GetOneByName(ctx context.Context, userID uuid.UUID, action models.PermissionAction) (*models.UserPermission, error) {
	row, err := s.queries.GetUserPermissionByName(ctx, db.GetUserPermissionByNameParams{
		UserID: models.UUIDPtrToDB(&userID),
		Name:   action.String(),
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

func (s *UserPermissionService) CreateByName(
	ctx context.Context,
	permAction models.PermissionAction,
	granter *models.User,
	granted bool,
	expiration *int64,
	grantedID uuid.UUID,
) error {
	dbPerm, err := s.pService.GetByName(ctx, permAction)
	if err != nil {
		return err
	}

	if granter != nil && granter.Role < dbPerm.RoleLevel {
		return errs.ErrPermissionDenied
	}

	var grantedBy *uuid.UUID
	if granter != nil {
		grantedBy = &granter.ID

	}

	perm := &models.UserPermission{
		Permission: *dbPerm,
		GrantedBy:  grantedBy,
		Granted:    granted,
		ExpiresAt:  expiration,
	}

	return s.Create(ctx, grantedID, perm)
}

func (s *UserPermissionService) Create(ctx context.Context, grantedID uuid.UUID, up *models.UserPermission) error {
	return s.queries.CreateUserPermission(ctx, *up.DBCreate(grantedID))
}

func (s *UserPermissionService) Update(ctx context.Context, uid, pid uuid.UUID, granted bool, expiresAt *int64) error {
	params := &models.UserPermission{
		Permission: models.Permission{ID: pid},
		Granted:    granted,
		ExpiresAt:  expiresAt,
	}

	return s.queries.UpdateUserPermission(ctx, *params.DBUpdate(uid))
}

func (s *UserPermissionService) Delete(ctx context.Context, uid, pid uuid.UUID) error {
	return s.queries.DeleteUserPermission(ctx, db.DeleteUserPermissionParams{
		PermissionID: models.UUIDPtrToDB(&pid),
		UserID:       models.UUIDPtrToDB(&uid),
	})
}
