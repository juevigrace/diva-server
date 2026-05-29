package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/util"
)

type AuthService struct {
	sService  *SessionService
	uService  *UserService
	uaService *UserActionsService
	uvService *UserVerificationService
}

func NewAuthService(
	uService *UserService,
	uaService *UserActionsService,
	uvService *UserVerificationService,
	sService *SessionService,
) *AuthService {
	return &AuthService{
		sService:  sService,
		uService:  uService,
		uaService: uaService,
		uvService: uvService,
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
		return nil, models.ErrInvalidCredentials
	}

	if !util.ValidatePassword(dto.Password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
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
		return nil, models.ErrSessionInvalid
	}
	updated, err := s.sService.Update(ctx, session)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *AuthService) ForgotPasswordConfirm(ctx context.Context, dto *dtos.ForgotPasswordConfirmDto) (*models.Session, error) {
	parsedID, err := uuid.Parse(dto.ActionID)
	if err != nil {
		return nil, err
	}

	dbUV, err := s.uvService.GetOneById(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	if !dbUV.Verified {
		return nil, models.ErrActionNotVerified
	}

	session, err := s.sService.CreateTemporal(ctx, dbUV.Action.UserID, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	// TODO: create user permission to update the password

	go func() {
		if err := s.uaService.Delete(ctx, dbUV.Action.ID); err != nil {
			log.Printf("action error: %v", err)
		}
	}()

	return session, nil
}

func (s *AuthService) SignOut(ctx context.Context, sessionID uuid.UUID) error {
	return s.sService.Close(ctx, sessionID)
}
