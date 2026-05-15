package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionRepository struct {
	queries *db.Queries
}

func NewSessionRepository(queries *db.Queries) *SessionRepository {
	return &SessionRepository{queries: queries}
}

func (r *SessionRepository) ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	rows, err := r.queries.ListSessionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	sessions := make([]*models.Session, len(rows))
	for i := range rows {
		sessions[i] = &models.Session{
			ID:           rows[i].ID.Bytes,
			User:         models.User{ID: rows[i].Userid.Bytes},
			AccessToken:  rows[i].Accesstoken,
			RefreshToken: rows[i].Refreshtoken,
			Device:       rows[i].Device,
			Status:       models.SessionStatusFromDB(rows[i].Status),
			Type:         models.SessionTypeFromDB(rows[i].Type),
			IpAddress:    rows[i].Ipaddress,
			UserAgent:    rows[i].Useragent,
			ExpiresAt:    rows[i].Expiresat.Time.UnixMilli(),
			CreatedAt:    rows[i].Createdat.Time.UnixMilli(),
			UpdatedAt:    rows[i].Updatedat.Time.UnixMilli(),
		}
	}
	return sessions, nil
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	row, err := r.queries.GetSessionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.Session{
		ID:           row.ID.Bytes,
		User:         models.User{ID: row.Userid.Bytes},
		AccessToken:  row.Accesstoken,
		RefreshToken: row.Refreshtoken,
		Device:       row.Device,
		Status:       models.SessionStatusFromDB(row.Status),
		Type:         models.SessionTypeFromDB(row.Type),
		IpAddress:    row.Ipaddress,
		UserAgent:    row.Useragent,
		ExpiresAt:    row.Expiresat.Time.UnixMilli(),
		CreatedAt:    row.Createdat.Time.UnixMilli(),
		UpdatedAt:    row.Updatedat.Time.UnixMilli(),
	}, nil
}

func (r *SessionRepository) GetByAccessToken(ctx context.Context, accessToken string) (*models.Session, error) {
	row, err := r.queries.GetSessionByAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return &models.Session{
		ID:           row.ID.Bytes,
		User:         models.User{ID: row.Userid.Bytes},
		AccessToken:  row.Accesstoken,
		RefreshToken: row.Refreshtoken,
		Device:       row.Device,
		Status:       models.SessionStatusFromDB(row.Status),
		Type:         models.SessionTypeFromDB(row.Type),
		IpAddress:    row.Ipaddress,
		UserAgent:    row.Useragent,
		ExpiresAt:    row.Expiresat.Time.UnixMilli(),
		CreatedAt:    row.Createdat.Time.UnixMilli(),
		UpdatedAt:    row.Updatedat.Time.UnixMilli(),
	}, nil
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	row, err := r.queries.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return &models.Session{
		ID:           row.ID.Bytes,
		User:         models.User{ID: row.Userid.Bytes},
		AccessToken:  row.Accesstoken,
		RefreshToken: row.Refreshtoken,
		Device:       row.Device,
		Status:       models.SessionStatusFromDB(row.Status),
		Type:         models.SessionTypeFromDB(row.Type),
		IpAddress:    row.Ipaddress,
		UserAgent:    row.Useragent,
		ExpiresAt:    row.Expiresat.Time.UnixMilli(),
		CreatedAt:    row.Createdat.Time.UnixMilli(),
		UpdatedAt:    row.Updatedat.Time.UnixMilli(),
	}, nil
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: session.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: session.User.ID, Valid: true},
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		Device:       session.Device,
		Status:       session.Status.ToDB(),
		Type:         session.Type.ToDB(),
		IpAddress:    session.IpAddress,
		UserAgent:    session.UserAgent,
		ExpiresAt:    pgtype.Timestamptz{Time: time.UnixMilli(session.ExpiresAt), Valid: true},
	})
}

func (r *SessionRepository) Update(ctx context.Context, session *models.Session) error {
	return r.queries.UpdateSession(ctx, db.UpdateSessionParams{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		Device:       session.Device,
		Status:       session.Status.ToDB(),
		Type:         session.Type.ToDB(),
		IpAddress:    session.IpAddress,
		UserAgent:    session.UserAgent,
		ExpiresAt:    pgtype.Timestamptz{Time: time.UnixMilli(session.ExpiresAt), Valid: true},
		ID:           pgtype.UUID{Bytes: session.ID, Valid: true},
	})
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, status models.SessionStatus, id uuid.UUID) error {
	return r.queries.UpdateSessionStatus(ctx, db.UpdateSessionStatusParams{
		Status: status.ToDB(),
		ID:     pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSession(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *SessionRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteSessionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	return r.queries.DeleteExpiredSessions(ctx)
}
