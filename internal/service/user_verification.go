package service

import (
	"context"
	"database/sql"
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
	sService  *SessionService
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
		sService:  sService,
		uaService: uaService,
	}
}

func (s *UserVerificationService) RequestVerification(
	ctx context.Context,
	sendTo *models.User,
	action models.Action,
) error {
	a, err := s.uaService.GetOneByName(ctx, sendTo.ID, action)
	if err != nil {
		return err
	}

	if err := s.GenerateAndSend(ctx, sendTo.Email, a); err != nil {
		return err
	}

	return nil
}

func (s *UserVerificationService) GenerateAndSend(ctx context.Context, email string, action *models.UserAction) error {
	exVerification, err := s.repo.GetByActionId(ctx, &action.ID)
	if err != nil {
		if ok := !errors.Is(sql.ErrNoRows, err); !ok {
			return err
		}
	}

	if exVerification != nil {
		if exVerification.ExpiresAt.Before(time.Now().UTC()) {
			if err := s.DeleteToken(ctx, exVerification.Token); err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	token, err := util.GenerateOTPCode()
	if err != nil {
		return err
	}

	params := &models.UserActionVerification{
		UserAction: action,
		Token:      token,
		ExpiresAt:  time.Now().UTC().Add(15 * time.Minute),
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return err
	}

	verification, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	if err := s.mail.SendVerificationEmail(ctx, email, verification); err != nil {
		if err := s.repo.DeleteByToken(ctx, token); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (s *UserVerificationService) Verify(ctx context.Context, token string) (*models.UserActionVerification, error) {
	record, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		if ok := errors.Is(pgx.ErrNoRows, err); ok {
			return nil, errors.New("token is not valid")
		} else {
			return nil, err
		}
	}

	if record.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("token expired")
	}

	return record, nil
}

func (s *UserVerificationService) Delete(ctx context.Context, uv *models.UserActionVerification) error {
	if err := s.uaService.Delete(ctx, uv.UserAction); err != nil {
		return err
	}

	return s.DeleteToken(ctx, uv.Token)
}

func (s *UserVerificationService) DeleteToken(ctx context.Context, token string) error {
	return s.repo.DeleteByToken(ctx, token)
}

func (s *UserVerificationService) HandlePasswordReset(ctx context.Context, userID *uuid.UUID, sessionData *dtos.SessionDataDto) (*responses.SessionResponse, error) {
	session, err := s.sService.Create(ctx, userID, sessionData)
	if err != nil {
		return nil, err
	}

	session.Type = models.SESSION_TEMPORAL
	if _, err := s.sService.Update(ctx, session); err != nil {
		return nil, err
	}

	session, err = s.sService.GetByID(ctx, &session.ID)
	if err != nil {
		return nil, err
	}

	return models.ToSessionResponse(session), nil
}

func (s *UserVerificationService) HandleVerifyUser(ctx context.Context, userID uuid.UUID) error {
	return s.uService.VerifyUser(ctx, &userID)
}
