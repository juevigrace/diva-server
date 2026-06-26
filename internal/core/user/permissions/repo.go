package permissions

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/core/permission"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/pkg/errs"
	"github.com/juevigrace/diva-server/storage"
)

type UserPermissionRepo struct {
	store storage.UserPermissionStore
	pRepo *permission.PermissionRepo
}

func NewUserPermissionRepo(store storage.UserPermissionStore, pRepo *permission.PermissionRepo) *UserPermissionRepo {
	return &UserPermissionRepo{
		pRepo: pRepo,
		store: store,
	}
}

func (s *UserPermissionRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPermission, error) {
	rows, err := s.store.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	perms := make([]models.UserPermission, len(rows))
	for i := range rows {
		perm, err := s.pRepo.GetByID(ctx, rows[i].PermissionID)
		if err != nil {
			continue
		}

		perms[i] = *models.UserPermissionFromDB(&rows[i], perm)
	}

	return perms, nil
}

func (s *UserPermissionRepo) GetOneByPermID(ctx context.Context, userID, permissionID uuid.UUID) (*models.UserPermission, error) {
	row, err := s.store.GetUserPermission(ctx, &storage.GetUserPermissionParams{
		PermissionID: permissionID,
		UserID:       userID,
	})
	if err != nil {
		return nil, err
	}

	dbPerm, err := s.pRepo.GetByID(ctx, row.PermissionID)
	if err != nil {
		return nil, err
	}

	return models.UserPermissionFromDB(row, dbPerm), nil
}

func (s *UserPermissionRepo) GetOneByName(ctx context.Context, userID uuid.UUID, action models.PermissionAction) (*models.UserPermission, error) {
	row, err := s.store.GetUserPermissionByName(ctx, &storage.GetUserPermissionByNameParams{
		UserID: userID,
		Name:   action.String(),
	})
	if err != nil {
		return nil, err
	}

	dbPerm, err := s.pRepo.GetByID(ctx, row.PermissionID)
	if err != nil {
		return nil, err
	}

	return models.UserPermissionFromDB(row, dbPerm), nil
}

func (s *UserPermissionRepo) CreateByName(
	ctx context.Context,
	permAction models.PermissionAction,
	granter *models.User,
	granted bool,
	expiration *int64,
	grantedID uuid.UUID,
) error {
	dbPerm, err := s.pRepo.GetByName(ctx, permAction)
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

func (s *UserPermissionRepo) Create(ctx context.Context, grantedID uuid.UUID, up *models.UserPermission) error {
	return s.store.CreateUserPermission(ctx, up.DBCreate(grantedID))
}

func (s *UserPermissionRepo) Update(ctx context.Context, uid, pid uuid.UUID, granted bool, expiresAt *int64) error {
	params := &models.UserPermission{
		Permission: models.Permission{ID: pid},
		Granted:    granted,
		ExpiresAt:  expiresAt,
	}

	return s.store.UpdateUserPermission(ctx, params.DBUpdate(uid))
}

func (s *UserPermissionRepo) Delete(ctx context.Context, uid, pid uuid.UUID) error {
	return s.store.DeleteUserPermission(ctx, &storage.DeleteUserPermissionParams{
		PermissionID: pid,
		UserID:       uid,
	})
}
