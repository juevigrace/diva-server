package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/pkg/jwt"
	"github.com/juevigrace/diva-server/storage"
)

type SessionRepo struct {
	store storage.SessionStore
}

func NewSessionRepo(store storage.SessionStore) *SessionRepo {
	return &SessionRepo{
		store: store,
	}
}

func (s *SessionRepo) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	rows, err := s.store.ListSessionsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessions := make([]*models.Session, len(rows))
	for i := range rows {
		sessions[i] = models.SessionFromDB(&rows[i])
	}
	return sessions, nil
}

func (s *SessionRepo) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	row, err := s.store.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return models.SessionFromDB(row), nil
}

func (s *SessionRepo) Create(ctx context.Context, userID uuid.UUID, sType models.SessionType, dto *dtos.SessionDataDto) (*models.Session, error) {
	sessionID := uuid.New()
	accessToken, err := jwt.CreateToken(sessionID, time.Duration(jwt.AccessExp)*time.Second)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken(sessionID, time.Duration(jwt.RefreshExp)*time.Second)
	if err != nil {
		return nil, err
	}

	accessExpiration, err := jwt.GetTokenExpiration(accessToken)
	if err != nil {
		return nil, err
	}
	refreshExpiration, err := jwt.GetTokenExpiration(refreshToken)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:              sessionID,
		User:            models.User{ID: userID},
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		Device:          dto.Device,
		IpAddress:       dto.IpAddress,
		UserAgent:       dto.UserAgent,
		Status:          models.SESSION_ACTIVE,
		Type:            sType,
		AccessExpiresAt:  accessExpiration.UnixMilli(),
		RefreshExpiresAt: refreshExpiration.UnixMilli(),
	}

	if err := s.store.CreateSession(ctx, session.DBCreate()); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, sessionID)
}

func (s *SessionRepo) CreateTemporal(ctx context.Context, userID uuid.UUID, dto *dtos.SessionDataDto) (*models.Session, error) {
	return s.Create(ctx, userID, models.SESSION_TEMPORAL, dto)
}

func (s *SessionRepo) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	accessToken, err := jwt.CreateToken(session.ID, time.Duration(jwt.AccessExp)*time.Second)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken(session.ID, time.Duration(jwt.RefreshExp)*time.Second)
	if err != nil {
		return nil, err
	}

	accessExpiration, err := jwt.GetTokenExpiration(accessToken)
	if err != nil {
		return nil, err
	}
	refreshExpiration, err := jwt.GetTokenExpiration(refreshToken)
	if err != nil {
		return nil, err
	}

	session.AccessToken = accessToken
	session.RefreshToken = refreshToken
	session.AccessExpiresAt = accessExpiration.UnixMilli()
	session.RefreshExpiresAt = refreshExpiration.UnixMilli()

	if err := s.store.UpdateSession(ctx, session.DBUpdate()); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, session.ID)
}

func (s *SessionRepo) UpdateStatus(ctx context.Context, status models.SessionStatus, sessionID uuid.UUID) error {
	return s.store.UpdateSessionStatus(ctx, &storage.UpdateSessionStatusParams{
		Status: status.ToDB(),
		ID:     sessionID,
	})

}

func (s *SessionRepo) Expire(ctx context.Context, sessionID uuid.UUID) error {
	return s.UpdateStatus(ctx, models.SESSION_EXPIRED, sessionID)
}

func (s *SessionRepo) Close(ctx context.Context, sessionID uuid.UUID) error {
	return s.UpdateStatus(ctx, models.SESSION_CLOSED, sessionID)
}

func (s *SessionRepo) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return s.store.DeleteSession(ctx, sessionID)
}

func (s *SessionRepo) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.store.DeleteSessionsByUser(ctx, userID)
}

func (s *SessionRepo) DeleteExpired(ctx context.Context) error {
	return s.store.DeleteExpiredSessions(ctx)
}

func (s *SessionRepo) CloseAllByUser(ctx context.Context, userID uuid.UUID) error {
	sessions, err := s.GetByUser(ctx, userID)
	if err != nil {
		return err
	}

	errs := make([]error, len(sessions))
	for i, session := range sessions {
		if err := s.Close(ctx, session.ID); err != nil {
			errs[i] = err
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
