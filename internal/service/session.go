package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/internal/util"
)

type SessionService struct {
	repo *repo.SessionRepository
}

func NewSessionService(repo *repo.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Create(ctx context.Context, userID *uuid.UUID, dto *dtos.SessionDataDto) (*models.Session, error) {
	sessionID := uuid.New()
	accessToken, err := util.CreateAccessToken(*userID, sessionID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.CreateRefreshToken(*userID, sessionID)
	if err != nil {
		return nil, err
	}

	expiration, err := util.GetTokenExpiration(refreshToken)
	if err != nil {
		return nil, err
	}
	expirationMillis := expiration.UnixMilli()

	session := &models.Session{
		ID:           sessionID,
		User:         models.User{ID: *userID},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Device:       dto.Device,
		IpAddress:    dto.IpAddress,
		UserAgent:    dto.UserAgent,
		Status:       models.SESSION_ACTIVE,
		ExpiresAt:    expirationMillis,
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, &sessionID)
}

func (s *SessionService) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	accessToken, err := util.CreateAccessToken(session.User.ID, session.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.CreateRefreshToken(session.User.ID, session.ID)
	if err != nil {
		return nil, err
	}

	expiration, err := util.GetTokenExpiration(refreshToken)
	if err != nil {
		return nil, err
	}
	expirationMillis := expiration.UnixMilli()

	var newSession models.Session = *session
	newSession.AccessToken = accessToken
	newSession.RefreshToken = refreshToken
	newSession.ExpiresAt = expirationMillis

	if err := s.repo.Update(ctx, &newSession); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, &session.ID)
}

func (s *SessionService) Close(ctx context.Context, sessionID *uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, models.SESSION_CLOSED, sessionID)
}

func (s *SessionService) Delete(ctx context.Context, sessionID *uuid.UUID) error {
	return s.repo.Delete(ctx, sessionID)
}

func (s *SessionService) GetByID(ctx context.Context, sessionID *uuid.UUID) (*models.Session, error) {
	return s.repo.GetByID(ctx, sessionID)
}

func (s *SessionService) GetByUser(ctx context.Context, userID *uuid.UUID) ([]*models.Session, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *SessionService) CloseAll(ctx context.Context, userID *uuid.UUID) error {
	sessions, err := s.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	errs := make([]error, 0)

	for _, session := range sessions {
		if err := s.Close(ctx, &session.ID); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
