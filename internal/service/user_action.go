package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserActionsService struct {
	repo *repo.UserActionsRepo
}

func NewUserActionsService(repo *repo.UserActionsRepo) *UserActionsService {
	return &UserActionsService{
		repo: repo,
	}
}

func (s *UserActionsService) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserAction, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *UserActionsService) GetOneByID(ctx context.Context, id uuid.UUID) (*models.UserAction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserActionsService) GetOneByName(ctx context.Context, userID uuid.UUID, action models.Action) (*models.UserAction, error) {
	return s.repo.GetByUserAndName(ctx, userID, action)
}

func (s *UserActionsService) Create(ctx context.Context, userID uuid.UUID, action models.Action) (*uuid.UUID, error) {
	id := uuid.New()

	if err := s.repo.Create(ctx, userID, &models.UserAction{
		ID:     id,
		Action: action,
	}); err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *UserActionsService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserActionsService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByUser(ctx, userID)
}
