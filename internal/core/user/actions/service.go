package actions

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/pkg/errs"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserActionsService struct {
	queries *db.Queries
}

func NewUserActionsService(queries *db.Queries) *UserActionsService {
	return &UserActionsService{
		queries: queries,
	}
}

func (s *UserActionsService) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.UserAction, error) {
	rows, err := s.queries.ListActionsByUser(ctx, models.UUIDPtrToDB(&userID))
	if err != nil {
		return nil, err
	}

	actions := make([]models.UserAction, len(rows))
	for i := range rows {
		actions[i] = *models.UserActionFromDB(&rows[i])
	}

	return actions, nil
}

func (s *UserActionsService) GetOneByID(ctx context.Context, id uuid.UUID) (*models.UserAction, error) {
	row, err := s.queries.GetUserActionByID(ctx, models.UUIDPtrToDB(&id))
	if err != nil {
		if ok := errors.Is(err, pgx.ErrNoRows); ok {
			return nil, errs.ErrActionNotFound
		} else {
			return nil, err
		}
	}

	return models.UserActionFromDB(&row), nil
}

func (s *UserActionsService) GetOneByName(ctx context.Context, userID uuid.UUID, action models.Action) (*models.UserAction, error) {
	row, err := s.queries.GetUserActionByUserAndName(ctx, db.GetUserActionByUserAndNameParams{
		UserID: models.UUIDPtrToDB(&userID),
		Name:   action.String(),
	})
	if err != nil {
		return nil, err
	}

	return models.UserActionFromDB(&row), nil
}

func (s *UserActionsService) Create(ctx context.Context, userID uuid.UUID, action models.Action) (*uuid.UUID, error) {
	params := &models.UserAction{
		ID:     uuid.New(),
		Name:   action,
		UserID: userID,
	}

	if err := s.queries.CreateUserAction(ctx, *params.DBCreate()); err != nil {
		return nil, err
	}

	return &params.ID, nil
}

func (s *UserActionsService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteUserAction(ctx, models.UUIDPtrToDB(&id))
}

func (s *UserActionsService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.queries.DeleteUserActionByUser(ctx, models.UUIDPtrToDB(&userID))
}
