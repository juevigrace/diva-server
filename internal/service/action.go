package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/repo"
)

type ActionService struct {
	repo *repo.ActionRepository
}

func NewActionService(repo *repo.ActionRepository) *ActionService {
	return &ActionService{repo: repo}
}

func (s *ActionService) Create(ctx context.Context, userAction *models.UserAction) error {
	return s.repo.Create(ctx, userAction)
}

func (s *ActionService) GetAll(ctx context.Context, userID uuid.UUID) ([]models.Action, error) {
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

func (s *ActionService) Delete(ctx context.Context, userAction *models.UserAction) error {
	return s.repo.Delete(ctx, userAction)
}
