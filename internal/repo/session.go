package repo

import (
	"context"

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

func (r *SessionRepository) Create(ctx context.Context, params *db.CreateSessionParams) error {
	return r.queries.CreateSession(ctx, *params)
}

func (r *SessionRepository) Update(ctx context.Context, params *db.UpdateSessionParams) error {
	return r.queries.UpdateSession(ctx, *params)
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, params *db.UpdateSessionStatusParams) error {
	return r.queries.UpdateSessionStatus(ctx, *params)
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
		BirthDate:    uRow.BirthDate.Time.Unix(),
		PhoneNumber:  uRow.PhoneNumber,
		Alias:        uRow.Alias,
		Avatar:       uRow.Avatar,
		Bio:          uRow.Bio,
		UserVerified: uRow.UserVerified,
		Role:         models.RoleFromDB(uRow.Role),
		CreatedAt:    uRow.CreatedAt.Time.Unix(),
		UpdatedAt:    uRow.UpdatedAt.Time.Unix(),
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
