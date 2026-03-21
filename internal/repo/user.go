package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.Count(ctx)
}

func (r *UserRepository) Create(ctx context.Context, params *db.CreateUserParams) (uuid.UUID, error) {
	if err := r.queries.CreateUser(ctx, *params); err != nil {
		return uuid.Nil, err
	}

	return params.ID.Bytes, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, params *db.UpdateProfileParams) error {
	return r.queries.UpdateProfile(ctx, *params)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, params *db.UpdatePasswordParams) error {
	return r.queries.UpdatePassword(ctx, *params)
}

func (r *UserRepository) UpdatePhoneNumber(ctx context.Context, id uuid.UUID, phone string) error {
	params := db.UpdatePhoneNumberParams{
		PhoneNumber: phone,
		ID:          pgtype.UUID{Bytes: id, Valid: true},
	}
	return r.queries.UpdatePhoneNumber(ctx, params)
}

func (r *UserRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) error {
	params := db.UpdateUsernameParams{
		Username: username,
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	}
	return r.queries.UpdateUsername(ctx, params)
}

func (r *UserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, dto *dtos.UserEmailDto) error {
	params := db.UpdateEmailParams{
		Email: dto.Email,
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	}
	return r.queries.UpdateEmail(ctx, params)
}

func (r *UserRepository) VerifyUser(ctx context.Context, userID *uuid.UUID) error {
	return r.queries.UpdateVerified(ctx, db.UpdateVerifiedParams{
		UserVerified: true,
		ID:           pgtype.UUID{Bytes: *userID, Valid: true},
	})
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	rows, err := r.queries.GetAllUsers(ctx, db.GetAllUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(rows))
	for i, row := range rows {
		users[i] = &models.User{
			ID:           row.ID.Bytes,
			Email:        row.Email,
			Username:     row.Username,
			PasswordHash: row.PasswordHash,
			BirthDate:    row.BirthDate.Time.UnixMilli(),
			PhoneNumber:  row.PhoneNumber,
			Alias:        row.Alias,
			Avatar:       row.Avatar,
			Bio:          row.Bio,
			UserVerified: row.UserVerified,
			Role:         models.RoleFromDB(row.Role),
			CreatedAt:    row.CreatedAt.Time.UnixMilli(),
			UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
			DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
		}
	}
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:           row.ID.Bytes,
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		BirthDate:    row.BirthDate.Time.UnixMilli(),
		PhoneNumber:  row.PhoneNumber,
		Alias:        row.Alias,
		Avatar:       row.Avatar,
		Bio:          row.Bio,
		UserVerified: row.UserVerified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
	}, nil
}

func (r *UserRepository) GetByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error) {
	row, err := r.queries.GetUserByUsernameOrEmail(ctx, db.GetUserByUsernameOrEmailParams{
		Username: username,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:           row.ID.Bytes,
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		BirthDate:    row.BirthDate.Time.UnixMilli(),
		PhoneNumber:  row.PhoneNumber,
		Alias:        row.Alias,
		Avatar:       row.Avatar,
		Bio:          row.Bio,
		UserVerified: row.UserVerified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
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
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		BirthDate:    row.BirthDate.Time.UnixMilli(),
		PhoneNumber:  row.PhoneNumber,
		Alias:        row.Alias,
		Avatar:       row.Avatar,
		Bio:          row.Bio,
		UserVerified: row.UserVerified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
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
		Email:        row.Email,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		BirthDate:    row.BirthDate.Time.UnixMilli(),
		PhoneNumber:  row.PhoneNumber,
		Alias:        row.Alias,
		Avatar:       row.Avatar,
		Bio:          row.Bio,
		UserVerified: row.UserVerified,
		Role:         models.RoleFromDB(row.Role),
		CreatedAt:    row.CreatedAt.Time.UnixMilli(),
		UpdatedAt:    row.UpdatedAt.Time.UnixMilli(),
		DeletedAt:    models.ToInt64Ptr(row.DeletedAt),
	}, nil
}
