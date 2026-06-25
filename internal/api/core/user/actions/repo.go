package actions

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/pkg/errs"
	"github.com/juevigrace/diva-server/storage"
)

type UserActionsRepo struct {
	store storage.UserActionStore
}

func NewUserActionsRepo(store storage.UserActionStore) *UserActionsRepo {
	return &UserActionsRepo{
		store: store,
	}
}

func (s *UserActionsRepo) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.UserAction, error) {
	rows, err := s.store.ListActionsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	actions := make([]models.UserAction, len(rows))
	for i := range rows {
		actions[i] = *models.UserActionFromDB(&rows[i])
	}

	return actions, nil
}

func (s *UserActionsRepo) GetOneByID(ctx context.Context, id uuid.UUID) (*models.UserAction, error) {
	row, err := s.store.GetUserActionByID(ctx, id)
	if err != nil {
		if ok := errors.Is(err, pgx.ErrNoRows); ok {
			return nil, errs.ErrActionNotFound
		} else {
			return nil, err
		}
	}

	return models.UserActionFromDB(row), nil
}

func (s *UserActionsRepo) GetOneByName(ctx context.Context, userID uuid.UUID, action models.Action) (*models.UserAction, error) {
	row, err := s.store.GetUserActionByUserAndName(ctx, &storage.GetUserActionByUserAndNameParams{
		UserID: userID,
		Name:   action.String(),
	})
	if err != nil {
		return nil, err
	}

	return models.UserActionFromDB(row), nil
}

func (s *UserActionsRepo) Create(ctx context.Context, userID uuid.UUID, action models.Action) (*uuid.UUID, error) {
	params := &models.UserAction{
		ID:     uuid.New(),
		Name:   action,
		UserID: userID,
	}

	if err := s.store.CreateUserAction(ctx, params.DBCreate()); err != nil {
		return nil, err
	}

	return &params.ID, nil
}

func (s *UserActionsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteUserAction(ctx, id)
}

func (s *UserActionsRepo) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.store.DeleteUserActionByUser(ctx, userID)
}
