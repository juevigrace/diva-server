package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionService struct {
	repo *repo.SessionRepository
}

func NewSessionService(repo *repo.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Create(ctx context.Context, userID uuid.UUID, dto *dtos.SessionDataDto) (*models.Session, error) {
	sessionID := uuid.New()
	accessToken, err := util.CreateAccessToken(userID, sessionID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.CreateRefreshToken(userID, sessionID)
	if err != nil {
		return nil, err
	}

	expiration, err := util.GetTokenExpiration(refreshToken)
	if err != nil {
		return nil, err
	}
	expirationMillis := expiration.UnixMilli()

	params := &db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: sessionID, Valid: true},
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Device:       dto.Device,
		Status:       models.SESSION_ACTIVE.ToDB(),
		IpAddress:    dto.IpAddress,
		UserAgent:    dto.UserAgent,
		ExpiresAt:    models.ToTimestamptzPtr(&expirationMillis),
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, sessionID)
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

	params := &db.UpdateSessionParams{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Device:       session.Device,
		Status:       models.SESSION_ACTIVE.ToDB(),
		IpAddress:    session.IpAddress,
		UserAgent:    session.UserAgent,
		ExpiresAt:    models.ToTimestamptzPtr(&expirationMillis),
		ID:           pgtype.UUID{Bytes: session.ID, Valid: true},
	}
	if err := s.repo.Update(ctx, params); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, session.ID)
}

func (s *SessionService) Close(ctx context.Context, sessionID *uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, &db.UpdateSessionStatusParams{
		Status: models.SESSION_CLOSED.ToDB(),
		ID:     pgtype.UUID{Bytes: *sessionID, Valid: true},
	})
}

func (s *SessionService) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.Delete(ctx, sessionID)
}

func (s *SessionService) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	return s.repo.GetByID(ctx, sessionID)
}
