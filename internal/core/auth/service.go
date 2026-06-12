package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/core"
	"github.com/juevigrace/diva-server/internal/core/permission"
	"github.com/juevigrace/diva-server/internal/core/session"
	"github.com/juevigrace/diva-server/internal/core/user"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/pkg/bcrypt"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type AuthService struct {
	pService   *permission.PermissionService
	uvProvider core.Provider[*models.UserActionVerification]
	sService   *session.SessionService
	uService   *user.UserService
}

func NewAuthService(
	pService *permission.PermissionService,
	sService *session.SessionService,
	uService *user.UserService,
) *AuthService {
	return &AuthService{
		pService: pService,
		sService: sService,
		uService: uService,
	}
}

func (s *AuthService) SignUp(ctx context.Context, dto *dtos.SignUpDto) (*models.Session, error) {
	userID, err := s.uService.Create(ctx, &dto.User)
	if err != nil {
		return nil, err
	}

	session, err := s.sService.Create(ctx, userID, models.SESSION_NORMAL, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) SignIn(ctx context.Context, dto *dtos.SignInDto) (*models.Session, error) {
	user, err := s.uService.GetByUsernameOrEmail(ctx, dto.Username)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	if !bcrypt.ValidatePassword(dto.Password, user.PasswordHash) {
		return nil, errs.ErrInvalidCredentials
	}

	session, err := s.sService.Create(ctx, user.ID, models.SESSION_NORMAL, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) Refresh(ctx context.Context, session *models.Session, dto *dtos.SessionDataDto) (*models.Session, error) {
	if session.Device != dto.Device || session.UserAgent != dto.UserAgent {
		if err := s.sService.Close(ctx, session.ID); err != nil {
			return nil, err
		}
		return nil, errs.ErrSessionInvalid
	}
	updated, err := s.sService.Update(ctx, session)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *AuthService) ForgotPasswordConfirm(ctx context.Context, actionID uuid.UUID, sd *dtos.SessionDataDto) (*models.Session, error) {
	dbUV, err := s.uvProvider.GetByID(ctx, actionID)
	if err != nil {
		return nil, err
	}

	if !dbUV.Verified {
		return nil, errs.ErrActionNotVerified
	}

	session, err := s.sService.CreateTemporal(ctx, dbUV.Action.UserID, sd)
	if err != nil {
		return nil, err
	}

	dbPerm, err := s.pService.GetByName(ctx, models.PERMISSION_USERS_PASSWORD_WRITE)
	if err != nil {
		return nil, err
	}

	exp := time.Now().UTC().Add(15 * time.Minute).UnixMilli()
	perm := &models.UserPermission{
		Permission: *dbPerm,
		GrantedBy:  nil,
		Granted:    true,
		ExpiresAt:  &exp,
	}

	if err := s.uService.upService.Create(ctx, dbUV.Action.UserID, perm); err != nil {
		return nil, err
	}

	if err := s.uService.uaService.Delete(ctx, dbUV.Action.ID); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) SignOut(ctx context.Context, sessionID uuid.UUID) error {
	return s.sService.Close(ctx, sessionID)
}
