package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/core/permission"
	"github.com/juevigrace/diva-server/internal/core/session"
	"github.com/juevigrace/diva-server/internal/core/user"
	"github.com/juevigrace/diva-server/internal/core/verification"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/pkg/bcrypt"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type AuthRepo struct {
	pRepo *permission.PermissionRepo
	sRepo *session.SessionRepo
	uRepo *user.UserRepo
	vRepo *verification.VerificationRepo
}

func NewAuthRepo(
	pRepo *permission.PermissionRepo,
	sRepo *session.SessionRepo,
	uRepo *user.UserRepo,
	vRepo *verification.VerificationRepo,
) *AuthRepo {
	return &AuthRepo{
		pRepo: pRepo,
		sRepo: sRepo,
		uRepo: uRepo,
		vRepo: vRepo,
	}
}

func (s *AuthRepo) SignUp(ctx context.Context, dto *dtos.SignUpDto) (*models.Session, error) {
	userID, err := s.uRepo.Create(ctx, &dto.User)
	if err != nil {
		return nil, err
	}

	session, err := s.sRepo.Create(ctx, userID, models.SESSION_NORMAL, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthRepo) SignIn(ctx context.Context, dto *dtos.SignInDto) (*models.Session, error) {
	user, err := s.uRepo.GetByUsernameOrEmail(ctx, dto.Username)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	if !bcrypt.ValidatePassword(dto.Password, user.PasswordHash) {
		return nil, errs.ErrInvalidCredentials
	}

	session, err := s.sRepo.Create(ctx, user.ID, models.SESSION_NORMAL, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthRepo) Refresh(ctx context.Context, session *models.Session, dto *dtos.SessionDataDto) (*models.Session, error) {
	if session.Device != dto.Device || session.UserAgent != dto.UserAgent {
		if err := s.sRepo.Close(ctx, session.ID); err != nil {
			return nil, err
		}
		return nil, errs.ErrSessionInvalid
	}
	updated, err := s.sRepo.Update(ctx, session)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *AuthRepo) ForgotPasswordConfirm(ctx context.Context, actionID uuid.UUID, sd *dtos.SessionDataDto) (*models.Session, error) {
	dbUV, err := s.vRepo.GetByID(ctx, actionID)
	if err != nil {
		return nil, err
	}

	if !dbUV.Verified {
		return nil, errs.ErrActionNotVerified
	}

	session, err := s.sRepo.CreateTemporal(ctx, dbUV.Action.UserID, sd)
	if err != nil {
		return nil, err
	}

	if err := s.uRepo.UpdatePasswordConfirm(ctx, dbUV.Action.ID, dbUV.Action.UserID); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthRepo) SignOut(ctx context.Context, sessionID uuid.UUID) error {
	return s.sRepo.Close(ctx, sessionID)
}
