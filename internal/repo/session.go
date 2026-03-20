package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/util"
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionRepository struct {
	queries *db.Queries
}

func NewSessionRepository(queries *db.Queries) *SessionRepository {
	return &SessionRepository{queries: queries}
}

func (r *SessionRepository) Create(ctx context.Context, sessionID, userID uuid.UUID, dto *dtos.SessionDataDto) error {
	accessToken, err := util.CreateAccessToken(userID, sessionID)
	if err != nil {
		return err
	}

	refreshToken, err := util.CreateRefreshToken(userID, sessionID)
	if err != nil {
		return err
	}

	expiration, err := util.GetTokenExpiration(refreshToken)
	if err != nil {
		return err
	}
	expirationMillis := expiration.UnixMilli()

	params := db.CreateSessionParams{
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
	return r.queries.CreateSession(ctx, params)
}

func (r *SessionRepository) Update(ctx context.Context, sessionID, userID uuid.UUID, dto *dtos.SessionDataDto) error {
	accessToken, err := util.CreateAccessToken(userID, sessionID)
	if err != nil {
		return err
	}

	refreshToken, err := util.CreateRefreshToken(userID, sessionID)
	if err != nil {
		return err
	}

	expiration, err := util.GetTokenExpiration(refreshToken)
	if err != nil {
		return err
	}
	expirationMillis := expiration.UnixMilli()

	params := db.UpdateSessionParams{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Device:       dto.Device,
		Status:       models.SESSION_ACTIVE.ToDB(),
		IpAddress:    dto.IpAddress,
		UserAgent:    dto.UserAgent,
		ExpiresAt:    models.ToTimestamptzPtr(&expirationMillis),
		ID:           pgtype.UUID{Bytes: sessionID, Valid: true},
	}
	return r.queries.UpdateSession(ctx, params)
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteSessionByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSession(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	row, err := r.queries.GetSessionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.Session{
		ID:           row.ID.Bytes,
		User:         row.UserID.Bytes,
		AccessToken:  row.AccessToken,
		RefreshToken: row.RefreshToken,
		Device:       row.Device,
		Status:       models.SessionStatusFromDB(row.Status),
		IpAddress:    row.IpAddress,
		UserAgent:    row.UserAgent,
		ExpiresAt:    row.ExpiresAt.Time.UnixMilli(),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
	}, nil
}
