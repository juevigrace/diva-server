package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/repo"
)

type UserActionsService struct {
	repo *repo.UserActionsRepository
}

func NewUserActionsService(repo *repo.UserActionsRepository) *UserActionsService {
	return &UserActionsService{repo: repo}
}

func (s *UserActionsService) Create(ctx context.Context, action models.Action, userID *uuid.UUID) (*uuid.UUID, error) {
	id := uuid.New()

	if err := s.repo.Create(ctx, &models.UserAction{
		ID:     id,
		Action: action,
		UserID: *userID,
	}); err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *UserActionsService) GetAll(ctx context.Context, userID uuid.UUID) ([]models.Action, error) {
	actions, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]models.Action, len(actions))
	for i, a := range actions {
		result[i] = models.ActionFromString(a.ActionName)
	}
	return result, nil
}

func (s *UserActionsService) GetOne(ctx context.Context, action models.Action, userID *uuid.UUID) (*models.UserAction, error) {
	return s.repo.GetOne(ctx, action, userID)
}

func (s *UserActionsService) Delete(ctx context.Context, userAction *models.UserAction) error {
	return s.repo.Delete(ctx, userAction)
}
