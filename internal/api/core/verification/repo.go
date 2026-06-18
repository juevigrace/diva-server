package verification

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/juevigrace/diva-server/internal/api/core/user"
	"github.com/juevigrace/diva-server/internal/api/core/user/actions"
	"github.com/juevigrace/diva-server/internal/api/core/user/permissions"
	"github.com/juevigrace/diva-server/pkg/mail"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/pkg/errs"
	"github.com/juevigrace/diva-server/pkg/otp"
	"github.com/juevigrace/diva-server/storage/db"
)

type VerificationRepo struct {
	mail    *mail.Client
	queries *db.Queries
	uRepo   *user.UserRepo
	uaRepo  *actions.UserActionsRepo
	upRepo  *permissions.UserPermissionRepo
	usRepo  *user.UserStateRepo
}

func NewVerificationRepo(
	mail *mail.Client,
	queries *db.Queries,
	uRepo *user.UserRepo,
	uaRepo *actions.UserActionsRepo,
	upRepo *permissions.UserPermissionRepo,
	usRepo *user.UserStateRepo,
) *VerificationRepo {
	return &VerificationRepo{
		mail:    mail,
		queries: queries,
		uRepo:   uRepo,
		uaRepo:  uaRepo,
		upRepo:  upRepo,
		usRepo:  usRepo,
	}
}

func (s *VerificationRepo) GetByID(ctx context.Context, actionID uuid.UUID) (*models.UserActionVerification, error) {
	dbAction, err := s.uaRepo.GetOneByID(ctx, actionID)
	if err != nil {
		return nil, err
	}

	row, err := s.queries.GetUserVerification(ctx, models.UUIDPtrToDB(&actionID))
	if err != nil {
		return nil, err
	}

	return &models.UserActionVerification{
		Action:    *dbAction,
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
		UsedAt:    &row.UsedAt.Time,
		Verified:  row.Verified,
	}, nil
}

func (s *VerificationRepo) RequestVerification(
	ctx context.Context,
	email string,
	action string,
) (*models.UserAction, error) {
	parsedAction := models.ActionFromString(action)
	if parsedAction == -1 {
		return nil, errs.ErrActionNotFound
	}

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.ID.Valid {
		return nil, errs.ErrIDRequired
	}

	parsedID, err := uuid.Parse(user.ID.String())
	if err != nil {
		return nil, err
	}

	dbAction, err := s.uaRepo.GetOneByName(ctx, parsedID, parsedAction)
	if err != nil {
		return nil, err
	}

	verification, err := s.Generate(ctx, dbAction)
	if err != nil {
		return nil, err
	}
	verification.Action = *dbAction

	if err := s.mail.Send(ctx, user.Email, "Email verification", VerificationEmail(verification)); err != nil {
		if err := s.Delete(ctx, verification.Action.ID); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &verification.Action, nil
}

func (s *VerificationRepo) Generate(
	ctx context.Context,
	action *models.UserAction,
) (*models.UserActionVerification, error) {
	exists, err := s.GetByID(ctx, action.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			token, err := otp.GenerateOTPCode()
			if err != nil {
				return nil, err
			}

			params := &models.UserActionVerification{
				Action:    *action,
				Token:     token,
				ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
			}

			if err := s.queries.CreateUserVerification(ctx, *params.DBCreate()); err != nil {
				return nil, err
			}

			return params, nil
		} else {
			return nil, err
		}
	}

	if exists.ExpiresAt.Before(time.Now().UTC()) {
		if err := s.Delete(ctx, exists.Action.ID); err != nil {
			return nil, err
		}
		return s.Generate(ctx, action)
	}

	return exists, nil
}

func (s *VerificationRepo) Verify(ctx context.Context, actionID uuid.UUID, token string) error {
	record, err := s.GetByID(ctx, actionID)
	if err != nil {
		return err
	}

	if record.ExpiresAt.Before(time.Now().UTC()) {
		return errs.ErrTokenExpired
	}

	if record.Token != token {
		return errs.ErrTokenInvalid
	}

	params := db.UpdateUserVerificationParams{
		Verified: true,
		ActionID: models.UUIDPtrToDB(&actionID),
	}

	if err := s.queries.UpdateUserVerification(ctx, params); err != nil {
		return err
	}

	va, err := s.GetByID(ctx, actionID)
	if err != nil {
		return err
	}

	return s.HandleVerified(ctx, va)
}

func (s *VerificationRepo) HandleVerified(ctx context.Context, va *models.UserActionVerification) error {
	if !va.Verified {
		return errs.ErrActionNotVerified
	}

	switch va.Action.Name {
	case models.ActionPasswordUpdate:
		return nil
	case models.ActionUserRestore:
		if err := s.uRepo.Restore(ctx, va.Action.UserID); err != nil {
			return err
		}

	case models.ActionUserVerification:
		if err := s.usRepo.UpdateVerified(ctx, true, va.Action.UserID); err != nil {
			return err
		}
	case models.ActionEmailUpdate, models.ActionUsernameUpdate, models.ActionPhoneUpdate:
		var permAction models.PermissionAction
		switch va.Action.Name {
		case models.ActionEmailUpdate:
			permAction = models.PERMISSION_USERS_EMAIL_WRITE
		case models.ActionUsernameUpdate:
			permAction = models.PERMISSION_USERS_USERNAME_WRITE
		case models.ActionPhoneUpdate:
			permAction = models.PERMISSION_USERS_PHONE_WRITE
		}

		dbPerm, err := s.upRepo.GetOneByName(ctx, va.Action.UserID, permAction)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		exp := time.Now().UTC().Add(15 * time.Minute).UnixMilli()
		if dbPerm == nil {
			if err := s.upRepo.CreateByName(ctx, permAction, nil, true, &exp, va.Action.UserID); err != nil {
				return err
			}
		} else if dbPerm.ExpiresAt != nil && time.UnixMilli(*dbPerm.ExpiresAt).Before(time.Now().UTC()) {
			if err := s.upRepo.Update(ctx, va.Action.UserID, dbPerm.Permission.ID, true, &exp); err != nil {
				return err
			}
		}
	}

	if err := s.uaRepo.Delete(ctx, va.Action.ID); err != nil {
		return err
	}

	return nil
}

func (s *VerificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteUserVerification(ctx, models.UUIDPtrToDB(&id))
}
