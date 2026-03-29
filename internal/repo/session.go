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

func (r *SessionRepository) Create(ctx context.Context, params *models.Session) error {
	return r.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: params.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: params.User.ID, Valid: true},
		AccessToken:  params.AccessToken,
		RefreshToken: params.RefreshToken,
		Device:       params.Device,
		Status:       params.Status.ToDB(),
		IpAddress:    params.IpAddress,
		UserAgent:    params.UserAgent,
		ExpiresAt:    models.ToTimestamptzPtr(&params.ExpiresAt),
	})
}

func (r *SessionRepository) Update(ctx context.Context, params *models.Session) error {
	return r.queries.UpdateSession(ctx, db.UpdateSessionParams{
		AccessToken:  params.AccessToken,
		RefreshToken: params.RefreshToken,
		Device:       params.Device,
		Status:       params.Status.ToDB(),
		IpAddress:    params.IpAddress,
		UserAgent:    params.UserAgent,
		ExpiresAt:    pgtype.Timestamptz{Time: time.UnixMilli(params.ExpiresAt), Valid: true},
		ID:           pgtype.UUID{Bytes: params.ID, Valid: true},
	})
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, status models.SessionStatus, id *uuid.UUID) error {
	return r.queries.UpdateSessionStatus(ctx, db.UpdateSessionStatusParams{
		Status: status.ToDB(),
		ID:     pgtype.UUID{Bytes: *id, Valid: true},
	})
}

func (r *SessionRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	return r.queries.DeleteSession(ctx, pgtype.UUID{Bytes: *id, Valid: true})
}

func (r *SessionRepository) GetByUser(ctx context.Context, userID *uuid.UUID) ([]*models.Session, error) {
	rows, err := r.queries.GetSessionsByUser(ctx, pgtype.UUID{Bytes: *userID, Valid: true})
	if err != nil {
		return nil, err
	}

	sessions := make([]*models.Session, len(rows))
	for i := range rows {
		sessions[i] = &models.Session{
			ID:           rows[i].ID.Bytes,
			AccessToken:  rows[i].AccessToken,
			RefreshToken: rows[i].RefreshToken,
			Device:       rows[i].Device,
			Status:       models.SessionStatusFromDB(rows[i].Status),
			IpAddress:    rows[i].IpAddress,
			UserAgent:    rows[i].UserAgent,
			ExpiresAt:    rows[i].ExpiresAt.Time.UnixMilli(),
			CreatedAt:    rows[i].CreatedAt.Time.UnixMilli(),
			UpdatedAt:    rows[i].UpdatedAt.Time.UnixMilli(),
		}
	}

	return sessions, nil
}

func (r *SessionRepository) GetByID(ctx context.Context, id *uuid.UUID) (*models.Session, error) {
	row, err := r.queries.GetSessionByID(ctx, pgtype.UUID{Bytes: *id, Valid: true})
	if err != nil {
		return nil, err
	}

	uRow, err := r.queries.GetUserByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}

	// TODO: find permissions and add them
	user := models.User{
		ID:           uRow.ID.Bytes,
		Email:        uRow.Email,
		Username:     uRow.Username,
		PasswordHash: uRow.PasswordHash,
		BirthDate:    uRow.BirthDate.Time.UnixMilli(),
		PhoneNumber:  uRow.PhoneNumber,
		Alias:        uRow.Alias,
		Avatar:       uRow.Avatar,
		Bio:          uRow.Bio,
		UserVerified: uRow.UserVerified,
		Role:         models.RoleFromDB(uRow.Role),
		CreatedAt:    uRow.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    uRow.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(uRow.DeletedAt),
	}

	return &models.Session{
		ID:           row.ID.Bytes,
		User:         user,
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
