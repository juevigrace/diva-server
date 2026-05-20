package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/util"
)

type AuthService struct {
	userService    *UserService
	uvService      *UserVerificationService
	sessionService *SessionService
}

func NewAuthService(
	userService *UserService,
	uvService *UserVerificationService,
	sessionService *SessionService,
) *AuthService {
	return &AuthService{
		userService:    userService,
		sessionService: sessionService,
		uvService:      uvService,
	}
}

func (s *AuthService) SignUp(ctx context.Context, dto *dtos.SignUpDto) (*responses.SessionResponse, error) {
	userID, err := s.userService.Create(ctx, &dto.User)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionService.Create(ctx, userID, models.SESSION_NORMAL, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return models.ToSessionResponse(session), nil
}

func (s *AuthService) SignIn(ctx context.Context, dto *dtos.SignInDto) (*responses.SessionResponse, error) {
	user, err := s.userService.GetByUsernameOrEmail(ctx, dto.Username)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if !util.ValidatePassword(dto.Password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	session, err := s.sessionService.Create(ctx, user.ID, models.SESSION_NORMAL, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return models.ToSessionResponse(session), nil
}

func (s *AuthService) SignOut(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionService.Close(ctx, sessionID)
}

func (s *AuthService) Refresh(ctx context.Context, session *models.Session, dto *dtos.SessionDataDto) (*responses.SessionResponse, error) {
	if session.Device != dto.Device || session.UserAgent != dto.UserAgent {
		if err := s.sessionService.Close(ctx, session.ID); err != nil {
			return nil, err
		}
		return nil, models.ErrSessionInvalid
	}
	updated, err := s.sessionService.Update(ctx, session)
	if err != nil {
		return nil, err
	}
	return models.ToSessionResponse(updated), nil
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

	session, err := s.sessionService.CreateTemporal(ctx, dbUV.Action.UserID, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	go func() {}()

	return session, nil
}

func (s *AuthService) ForgotPasswordUpdate(
	ctx context.Context,
	session *models.Session,
	dto *dtos.UpdatePasswordDto,
) error {
	if err := s.userService.UpdatePassword(ctx, session, dto.NewPassword); err != nil {
		return err
	}

	if err := s.sessionService.Delete(ctx, session.ID); err != nil {
		return err
	}

	if err := s.sessionService.CloseAllByUser(ctx, session.User.ID); err != nil {
		return err
	}

	return nil
}
