package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/internal/util"
)

type UserVerificationService struct {
	mail      *mail.Client
	repo      *repo.UserVerificationRepo
	uaService *UserActionsService
}

func NewVerificationService(
	mail *mail.Client,
	repo *repo.UserVerificationRepo,
	sService *SessionService,
	uaService *UserActionsService,
) *UserVerificationService {
	return &UserVerificationService{
		mail:      mail,
		repo:      repo,
		uaService: uaService,
	}
}

func (s *UserVerificationService) GetOneById(ctx context.Context, actionID uuid.UUID) (*models.UserActionVerification, error) {
	dbAction, err := s.uaService.GetOneByID(ctx, actionID)
	if err != nil {
		return nil, err
	}

	dbUV, err := s.repo.GetByID(ctx, actionID)
	if err != nil {
		return nil, err
	}

	dbUV.Action = *dbAction

	return dbUV, nil
}

func (s *UserVerificationService) RequestVerification(
	ctx context.Context,
	user *models.User,
	dto *dtos.RequestActionVerificationDto,
) (*responses.UserActionResponse, error) {
	parsedAction := models.ActionFromString(dto.Action)
	if parsedAction == -1 {
		return nil, models.ErrActionNotFound
	}

	verification, err := s.Generate(ctx, user.ID, parsedAction)
	if err != nil {
		return nil, err
	}

	if err := s.mail.SendVerificationEmail(ctx, user.Email, verification); err != nil {
		if err := s.Delete(ctx, verification.Action.ID); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &responses.UserActionResponse{
		ID:         verification.Action.ID.String(),
		ActionName: verification.Action.Name.String(),
	}, nil
}

func (s *UserVerificationService) Generate(
	ctx context.Context,
	userID uuid.UUID,
	action models.Action,
) (*models.UserActionVerification, error) {
	dbAction, err := s.uaService.GetOneByName(ctx, userID, action)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.GetByID(ctx, dbAction.ID)
	if err != nil {
		if ok := errors.Is(pgx.ErrNoRows, err); ok {
			token, err := util.GenerateOTPCode()
			if err != nil {
				return nil, err
			}

			params := &models.UserActionVerification{
				Action:    *dbAction,
				Token:     token,
				ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
			}

			if err := s.repo.Create(ctx, params); err != nil {
				return nil, err
			}

			return s.Generate(ctx, userID, action)
		} else {
			return nil, err
		}
	}

	if exists.ExpiresAt.Before(time.Now().UTC()) {
		if err := s.Delete(ctx, exists.Action.ID); err != nil {
			return nil, err
		}
		return s.Generate(ctx, userID, action)
	}

	return exists, nil
}

func (s *UserVerificationService) Verify(ctx context.Context, actionID uuid.UUID, token string) (*models.UserAction, error) {
	dbAction, err := s.uaService.GetOneByID(ctx, actionID)
	if err != nil {
		return nil, err
	}

	record, err := s.repo.GetByID(ctx, dbAction.ID)
	if err != nil {
		if ok := errors.Is(pgx.ErrNoRows, err); ok {
			return nil, models.ErrActionNotFound
		} else {
			return nil, err
		}
	}

	if record.ExpiresAt.Before(time.Now().UTC()) {
		return nil, models.ErrTokenExpired
	}

	if record.Token != token {
		return nil, models.ErrTokenInvalid
	}

	if err := s.repo.Verify(ctx, dbAction.ID); err != nil {
		return nil, err
	}

	record.Action = *dbAction

	return &record.Action, nil
}

func (s *UserVerificationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
