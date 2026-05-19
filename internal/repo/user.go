package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserRepo struct {
	queries *db.Queries
}

func NewUserRepo(queries *db.Queries) *UserRepo {
	return &UserRepo{queries: queries}
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *UserRepo) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	rows, err := r.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(rows))
	for i := range rows {
		users[i] = &models.User{
			ID:           rows[i].ID.Bytes,
			Username:     rows[i].Username,
			Email:        rows[i].Email,
			PhoneNumber:  rows[i].PhoneNumber,
			PasswordHash: rows[i].PasswordHash,
			Verified:     rows[i].Verified,
			Role:         models.RoleFromDB(rows[i].Role),
			CreatedAt:    rows[i].CreatedAt.Time.UnixMilli(),
			UpdatedAt:    rows[i].UpdatedAt.Time.UnixMilli(),
			DeletedAt:    models.ToInt64Ptr(rows[i].DeletedAt),
		}
	}
	return users, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &models.User{
		ID:           row.ID.Bytes,
		Username:     row.Username,
		Email:        row.Email,
		PhoneNumber:  row.PhoneNumber,
		PasswordHash: row.PasswordHash,
		Verified:     row.Verified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
	}, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &models.User{
		ID:           row.ID.Bytes,
		Username:     row.Username,
		Email:        row.Email,
		PhoneNumber:  row.PhoneNumber,
		PasswordHash: row.PasswordHash,
		Verified:     row.Verified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
	}, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &models.User{
		ID:           row.ID.Bytes,
		Username:     row.Username,
		Email:        row.Email,
		PhoneNumber:  row.PhoneNumber,
		PasswordHash: row.PasswordHash,
		Verified:     row.Verified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
	}, nil
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           pgtype.UUID{Bytes: user.ID, Valid: true},
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Verified:     user.Verified,
		Role:         user.Role.ToDB(),
	})
}

func (r *UserRepo) UpdateUsername(ctx context.Context, username string, id uuid.UUID) error {
	return r.queries.UpdateUsername(ctx, db.UpdateUsernameParams{
		Username: username,
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepo) UpdateEmail(ctx context.Context, email string, id uuid.UUID) error {
	return r.queries.UpdateEmail(ctx, db.UpdateEmailParams{
		Email: email,
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepo) UpdatePassword(ctx context.Context, passwordHash string, id uuid.UUID) error {
	return r.queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		PasswordHash: passwordHash,
		ID:           pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepo) UpdateVerified(ctx context.Context, verified bool, id uuid.UUID) error {
	return r.queries.UpdateVerified(ctx, db.UpdateVerifiedParams{
		Verified: verified,
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepo) UpdateRole(ctx context.Context, role models.Role, id uuid.UUID) error {
	return r.queries.UpdateRole(ctx, db.UpdateRoleParams{
		Role: role.ToDB(),
		ID:   pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepo) UpdatePhoneNumber(ctx context.Context, phoneNumber string, userID uuid.UUID) error {
	return r.queries.UpdatePhoneNumber(ctx, db.UpdatePhoneNumberParams{
		PhoneNumber: phoneNumber,
		ID:          pgtype.UUID{Bytes: userID, Valid: true},
	})
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserRepo) Restore(ctx context.Context, id uuid.UUID) error {
	return r.queries.RestoreUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
