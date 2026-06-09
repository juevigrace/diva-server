package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserVerificationService struct {
	mail      *mail.Client
	queries   *db.Queries
	uaService *UserActionsService
}

func NewVerificationService(
	mail *mail.Client,
	queries *db.Queries,
	uaService *UserActionsService,
) *UserVerificationService {
	return &UserVerificationService{
		mail:      mail,
		queries:   queries,
		uaService: uaService,
	}
}

func (s *UserVerificationService) GetOneById(ctx context.Context, actionID uuid.UUID) (*models.UserActionVerification, error) {
	dbAction, err := s.uaService.GetOneByID(ctx, actionID)
	if err != nil {
		return nil, err
	}

	row, err := s.queries.GetUserVerification(ctx, models.UUIDPtrToDB(&actionID))
	if err != nil {
		return nil, err
	}

	return &models.UserActionVerification{
		Action:    *dbAction,
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
		UsedAt:    &row.UsedAt.Time,
		Verified:  row.Verified,
	}, nil
}

func (s *UserVerificationService) RequestVerification(
	ctx context.Context,
	user *models.User,
	action models.Action,
) (*models.UserAction, error) {
	dbAction, err := s.uaService.GetOneByName(ctx, user.ID, action)
	if err != nil {
		return nil, err
	}

	verification, err := s.Generate(ctx, dbAction)
	if err != nil {
		return nil, err
	}
	verification.Action = *dbAction

	if err := s.mail.SendVerificationEmail(ctx, user.Email, verification); err != nil {
		if err := s.Delete(ctx, verification.Action.ID); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &verification.Action, nil
}

func (s *UserVerificationService) Generate(
	ctx context.Context,
	action *models.UserAction,
) (*models.UserActionVerification, error) {
	exists, err := s.GetOneById(ctx, action.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			token, err := util.GenerateOTPCode()
			if err != nil {
				return nil, err
			}

			params := &models.UserActionVerification{
				Action:    *action,
				Token:     token,
				ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
			}

			if err := s.queries.CreateUserVerification(ctx, *params.DBCreate()); err != nil {
				return nil, err
			}

			return params, nil
		} else {
			return nil, err
		}
	}

	if exists.ExpiresAt.Before(time.Now().UTC()) {
		if err := s.Delete(ctx, exists.Action.ID); err != nil {
			return nil, err
		}
		return s.Generate(ctx, action)
	}

	return exists, nil
}

func (s *UserVerificationService) Verify(ctx context.Context, actionID uuid.UUID, token string) (*models.UserActionVerification, error) {
	record, err := s.GetOneById(ctx, actionID)
	if err != nil {
		return nil, err
	}

	if record.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errs.ErrTokenExpired
	}

	if record.Token != token {
		return nil, errs.ErrTokenInvalid
	}

	params := db.UpdateUserVerificationParams{
		Verified: true,
		ActionID: models.UUIDPtrToDB(&actionID),
	}

	if err := s.queries.UpdateUserVerification(ctx, params); err != nil {
		return nil, err
	}

	return s.GetOneById(ctx, actionID)
}

func (s *UserVerificationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteUserVerification(ctx, models.UUIDPtrToDB(&id))
}
