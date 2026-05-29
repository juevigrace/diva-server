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
	repo     *repo.UserPermsRepo
	pService *PermissionService
}

func NewUserPermissionService(repo *repo.UserPermsRepo, pService *PermissionService) *UserPermissionService {
	return &UserPermissionService{
		repo:     repo,
		pService: pService,
	}
}

func (s *UserPermissionService) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPermission, error) {
	dbUPerm, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, up := range dbUPerm {
		dbPerm, err := s.pService.GetByID(ctx, up.Permission.ID)
		if err != nil {
			continue
		}
		up.Permission = models.Permission{
			ID:          dbPerm.ID,
			Name:        dbPerm.Name,
			Description: dbPerm.Description,
			Action:      dbPerm.Action,
			RoleLevel:   dbPerm.RoleLevel,
			CreatedAt:   dbPerm.CreatedAt,
			UpdatedAt:   dbPerm.UpdatedAt,
			DeletedAt:   dbPerm.DeletedAt,
		}
	}

	return dbUPerm, nil
}

func (s *UserPermissionService) GetOneByUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) (*models.UserPermission, error) {
	dbUPerm, err := s.repo.GetOneByUser(ctx, userID, permissionID)
	if err != nil {
		return nil, err
	}

	dbPerm, err := s.pService.GetByID(ctx, dbUPerm.Permission.ID)
	if err != nil {
		return nil, err
	}

	dbUPerm.Permission = models.Permission{
		ID:          dbPerm.ID,
		Name:        dbPerm.Name,
		Description: dbPerm.Description,
		Action:      dbPerm.Action,
		RoleLevel:   dbPerm.RoleLevel,
		CreatedAt:   dbPerm.CreatedAt,
		UpdatedAt:   dbPerm.UpdatedAt,
		DeletedAt:   dbPerm.DeletedAt,
	}

	return dbUPerm, nil
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

func (s *UserPermissionService) Delete(ctx context.Context, uid, pid uuid.UUID) error {
	return s.repo.Delete(ctx, uid, pid)
}

func (s *UserPermissionService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByUser(ctx, userID)
}
