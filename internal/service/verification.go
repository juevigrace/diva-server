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

type VerificationService struct {
	mail      *mail.Client
	repo      *repo.VerificationRepository
	sService  *SessionService
	uService  *UserService
	uaService *UserActionsService
}

func NewVerificationService(
	mail *mail.Client,
	repo *repo.VerificationRepository,
	sService *SessionService,
	uService *UserService,
	uaService *UserActionsService,
) *VerificationService {
	return &VerificationService{
		mail:      mail,
		repo:      repo,
		sService:  sService,
		uService:  uService,
		uaService: uaService,
	}
}

func (s *VerificationService) RequestVerification(ctx context.Context, dto *dtos.RequestVerificationDto) error {
	parsedAction := models.ActionFromString(dto.Action)
	if parsedAction == -1 {
		// TODO: create proper error
		return errors.New("action doesn't exists")
	}

	u, err := s.uService.GetByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}

	var action *models.UserAction
	a, err := s.uaService.GetOne(ctx, parsedAction, &u.ID)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			_, err := s.uaService.Create(ctx, parsedAction, &u.ID)
			if err != nil {
				return err
			}
			action, err = s.uaService.GetOne(ctx, parsedAction, &u.ID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		action = a
	}

	if err := s.GenerateAndSend(ctx, u.Email, &models.UserAction{
		ID:     action.ID,
		Action: action.Action,
		UserID: action.UserID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *VerificationService) GenerateAndSend(ctx context.Context, email string, action *models.UserAction) error {
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

	params := &models.UserVerification{
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

func (s *VerificationService) Verify(ctx context.Context, token string) (*models.UserVerification, error) {
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

func (s *VerificationService) Delete(ctx context.Context, uv *models.UserVerification) error {
	if err := s.uaService.Delete(ctx, uv.UserAction); err != nil {
		return err
	}

	return s.DeleteToken(ctx, uv.Token)
}

func (s *VerificationService) DeleteToken(ctx context.Context, token string) error {
	return s.repo.DeleteByToken(ctx, token)
}

func (s *VerificationService) HandlePasswordReset(ctx context.Context, userID *uuid.UUID, sessionData *dtos.SessionDataDto) (*responses.SessionResponse, error) {
	session, err := s.sService.Create(ctx, userID, sessionData)
	if err != nil {
		return nil, err
	}

	return models.ToSessionResponse(session), nil
}

func (s *VerificationService) HandleVerifyUser(ctx context.Context, userID uuid.UUID) error {
	return s.uService.VerifyUser(ctx, &userID)
}
