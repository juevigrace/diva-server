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

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           pgtype.UUID{Bytes: user.ID, Valid: true},
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Alias:        user.Alias,
	})
}

func (r *UserRepository) CreateBatch(ctx context.Context, params []*models.User) error {
	for _, p := range params {
		if err := r.Create(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, params *models.User) error {
	return r.queries.UpdateProfile(ctx, db.UpdateProfileParams{
		Alias:  params.Alias,
		Avatar: params.Avatar,
		Bio:    params.Bio,
		ID:     pgtype.UUID{Bytes: params.ID, Valid: true},
	})
}

func (r *UserRepository) UpdatePassword(ctx context.Context, hash string, id *uuid.UUID) error {
	return r.queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		PasswordHash: hash,
		ID:           pgtype.UUID{Bytes: *id, Valid: true},
	})
}

func (r *UserRepository) UpdatePhoneNumber(ctx context.Context, phone string, id *uuid.UUID) error {
	params := db.UpdatePhoneNumberParams{
		PhoneNumber: phone,
		ID:          pgtype.UUID{Bytes: *id, Valid: true},
	}
	return r.queries.UpdatePhoneNumber(ctx, params)
}

func (r *UserRepository) UpdateUsername(ctx context.Context, username string, id *uuid.UUID) error {
	params := db.UpdateUsernameParams{
		Username: username,
		ID:       pgtype.UUID{Bytes: *id, Valid: true},
	}
	return r.queries.UpdateUsername(ctx, params)
}

func (r *UserRepository) UpdateEmail(ctx context.Context, dto *dtos.UpdateEmailDto, id *uuid.UUID) error {
	params := db.UpdateEmailParams{
		Email: dto.Email,
		ID:    pgtype.UUID{Bytes: *id, Valid: true},
	}
	return r.queries.UpdateEmail(ctx, params)
}

func (r *UserRepository) VerifyUser(ctx context.Context, userID *uuid.UUID) error {
	return r.queries.UpdateVerified(ctx, db.UpdateVerifiedParams{
		UserVerified: true,
		ID:           pgtype.UUID{Bytes: *userID, Valid: true},
	})
}

func (r *UserRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	return r.queries.DeleteUser(ctx, pgtype.UUID{Bytes: *id, Valid: true})
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

func (r *UserRepository) GetByID(ctx context.Context, id *uuid.UUID) (*models.User, error) {
	row, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: *id, Valid: true})
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
