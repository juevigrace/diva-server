package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionService struct {
	queries  *db.Queries
	uService *UserService
}

func NewSessionService(queries *db.Queries, uService *UserService) *SessionService {
	return &SessionService{
		queries:  queries,
		uService: uService,
	}
}

func (s *SessionService) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	rows, err := s.queries.ListSessionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	sessions := make([]*models.Session, len(rows))
	for i := range rows {
		sessions[i] = models.SessionFromDB(&rows[i])
	}
	return sessions, nil
}

func (s *SessionService) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	row, err := s.queries.GetSessionByID(ctx, models.UUIDPtrToDB(&sessionID))
	if err != nil {
		return nil, err
	}

	dbSession := models.SessionFromDB(&row)

	dbUser, err := s.uService.GetByID(ctx, dbSession.User.ID)
	if err != nil {
		return nil, err
	}

	dbSession.User = *dbUser

	return dbSession, nil
}

func (s *SessionService) Create(ctx context.Context, userID uuid.UUID, sType models.SessionType, dto *dtos.SessionDataDto) (*models.Session, error) {
	sessionID := uuid.New()
	accessToken, err := util.CreateAccessToken(sessionID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.CreateRefreshToken(sessionID)
	if err != nil {
		return nil, err
	}

	expiration, err := util.GetTokenExpiration(refreshToken)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:           sessionID,
		User:         models.User{ID: userID},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Device:       dto.Device,
		IpAddress:    dto.IpAddress,
		UserAgent:    dto.UserAgent,
		Status:       models.SESSION_ACTIVE,
		Type:         sType,
		ExpiresAt:    expiration.UnixMilli(),
	}

	if err := s.queries.CreateSession(ctx, *session.DBCreate()); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, sessionID)
}

func (s *SessionService) CreateTemporal(ctx context.Context, userID uuid.UUID, dto *dtos.SessionDataDto) (*models.Session, error) {
	return s.Create(ctx, userID, models.SESSION_TEMPORAL, dto)
}

func (s *SessionService) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	accessToken, err := util.CreateAccessToken(session.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.CreateRefreshToken(session.ID)
	if err != nil {
		return nil, err
	}

	expiration, err := util.GetTokenExpiration(refreshToken)
	if err != nil {
		return nil, err
	}

	session.AccessToken = accessToken
	session.RefreshToken = refreshToken
	session.ExpiresAt = expiration.UnixMilli()

	if err := s.queries.UpdateSession(ctx, *session.DBUpdate()); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, session.ID)
}

func (s *SessionService) UpdateStatus(ctx context.Context, status models.SessionStatus, sessionID uuid.UUID) error {
	return s.queries.UpdateSessionStatus(ctx, db.UpdateSessionStatusParams{
		Status: status.ToDB(),
		ID:     models.UUIDPtrToDB(&sessionID),
	})

}

func (s *SessionService) Expire(ctx context.Context, sessionID uuid.UUID) error {
	return s.UpdateStatus(ctx, models.SESSION_EXPIRED, sessionID)
}

func (s *SessionService) Close(ctx context.Context, sessionID uuid.UUID) error {
	return s.UpdateStatus(ctx, models.SESSION_CLOSED, sessionID)
}

func (s *SessionService) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return s.queries.DeleteSession(ctx, models.UUIDPtrToDB(&sessionID))
}

func (s *SessionService) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return s.queries.DeleteSessionsByUser(ctx, models.UUIDPtrToDB(&userID))
}

func (s *SessionService) DeleteExpired(ctx context.Context) error {
	return s.queries.DeleteExpiredSessions(ctx)
}

func (s *SessionService) CloseAllByUser(ctx context.Context, userID uuid.UUID) error {
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
