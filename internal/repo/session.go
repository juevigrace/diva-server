package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionRepo struct {
	queries *db.Queries
}

func NewSessionRepo(queries *db.Queries) *SessionRepo {
	return &SessionRepo{queries: queries}
}

func (r *SessionRepo) ListSessionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	rows, err := r.queries.ListSessionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	sessions := make([]*models.Session, len(rows))
	for i := range rows {
		sessions[i] = &models.Session{
			ID:           rows[i].ID.Bytes,
			User:         models.User{ID: rows[i].UserID.Bytes},
			AccessToken:  rows[i].AccessToken,
			RefreshToken: rows[i].RefreshToken,
			Device:       rows[i].Device,
			Status:       models.SessionStatusFromDB(rows[i].Status),
			Type:         models.SessionTypeFromDB(rows[i].Type),
			IpAddress:    rows[i].IpAddress,
			UserAgent:    rows[i].UserAgent,
			ExpiresAt:    rows[i].ExpiresAt.Time.UnixMilli(),
			CreatedAt:    rows[i].CreatedAt.Time.UnixMilli(),
			UpdatedAt:    rows[i].UpdatedAt.Time.UnixMilli(),
		}
	}
	return sessions, nil
}

func (r *SessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	row, err := r.queries.GetSessionByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.Session{
		ID:           row.ID.Bytes,
		User:         models.User{ID: row.UserID.Bytes},
		AccessToken:  row.AccessToken,
		RefreshToken: row.RefreshToken,
		Device:       row.Device,
		Status:       models.SessionStatusFromDB(row.Status),
		Type:         models.SessionTypeFromDB(row.Type),
		IpAddress:    row.IpAddress,
		UserAgent:    row.UserAgent,
		ExpiresAt:    row.ExpiresAt.Time.UnixMilli(),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
	}, nil
}

func (r *SessionRepo) Create(ctx context.Context, session *models.Session) error {
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

func (r *SessionRepo) Update(ctx context.Context, session *models.Session) error {
	return r.queries.UpdateSession(ctx, db.UpdateSessionParams{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		IpAddress:    session.IpAddress,
		ExpiresAt:    pgtype.Timestamptz{Time: time.UnixMilli(session.ExpiresAt), Valid: true},
		ID:           pgtype.UUID{Bytes: session.ID, Valid: true},
	})
}

func (r *SessionRepo) UpdateStatus(ctx context.Context, status models.SessionStatus, id uuid.UUID) error {
	return r.queries.UpdateSessionStatus(ctx, db.UpdateSessionStatusParams{
		Status: status.ToDB(),
		ID:     pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *SessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSession(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *SessionRepo) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteSessionsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (r *SessionRepo) DeleteExpired(ctx context.Context) error {
	return r.queries.DeleteExpiredSessions(ctx)
}
