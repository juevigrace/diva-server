package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(rows))
	for i := range rows {
		users[i] = &models.User{
			ID:           rows[i].ID.Bytes,
			Username:     rows[i].Username,
			Email:        rows[i].Email,
			PasswordHash: rows[i].Passwordhash,
			Verified:     rows[i].Verified,
			Role:         models.RoleFromDB(rows[i].Role),
			CreatedAt:    rows[i].Createdat.Time.UnixMilli(),
			UpdatedAt:    rows[i].Updatedat.Time.UnixMilli(),
			DeletedAt:    models.ToInt64Ptr(rows[i].Deletedat),
		}
	}
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
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
		PasswordHash: row.Passwordhash,
		Verified:     row.Verified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.Createdat.Time.UnixMilli(),
		UpdatedAt:    row.Updatedat.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.Deletedat),
	}, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
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
		PasswordHash: row.Passwordhash,
		Verified:     row.Verified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.Createdat.Time.UnixMilli(),
		UpdatedAt:    row.Updatedat.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.Deletedat),
	}, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
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
		PasswordHash: row.Passwordhash,
		Verified:     row.Verified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.Createdat.Time.UnixMilli(),
		UpdatedAt:    row.Updatedat.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.Deletedat),
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           pgtype.UUID{Bytes: user.ID, Valid: true},
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Verified:     user.Verified,
		Role:         user.Role.ToDB(),
		CreatedAt:    pgtype.Timestamptz{Time: time.UnixMilli(user.CreatedAt), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.UnixMilli(user.UpdatedAt), Valid: true},
	})
}

func (r *UserRepository) UpdateUsername(ctx context.Context, username string, id uuid.UUID) error {
	return r.queries.UpdateUsername(ctx, db.UpdateUsernameParams{
		Username: username,
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepository) UpdateEmail(ctx context.Context, email string, id uuid.UUID) error {
	return r.queries.UpdateEmail(ctx, db.UpdateEmailParams{
		Email: email,
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepository) UpdatePassword(ctx context.Context, passwordHash string, id uuid.UUID) error {
	return r.queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		PasswordHash: passwordHash,
		ID:           pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepository) UpdateVerified(ctx context.Context, verified bool, id uuid.UUID) error {
	return r.queries.UpdateVerified(ctx, db.UpdateVerifiedParams{
		Verified: verified,
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepository) UpdateRole(ctx context.Context, role models.Role, id uuid.UUID) error {
	return r.queries.UpdateRole(ctx, db.UpdateRoleParams{
		Role: role.ToDB(),
		ID:   pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return r.queries.RestoreUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
