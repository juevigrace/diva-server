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
	userService         *UserService
	sessionService      *SessionService
	verificationService *VerificationService
	actionService       *UserActionsService
}

func NewAuthService(
	userService *UserService,
	sessionService *SessionService,
	verificationService *VerificationService,
	actionService *UserActionsService,
) *AuthService {
	return &AuthService{
		userService:         userService,
		sessionService:      sessionService,
		verificationService: verificationService,
		actionService:       actionService,
	}
}

func (s *AuthService) SignUp(ctx context.Context, dto *dtos.SignUpDto) (*responses.SessionResponse, error) {
	userID, err := s.userService.Create(ctx, &dto.User)
	if err != nil {
		return nil, err
	}

	if err := s.verificationService.GenerateAndSend(ctx, userID, dto.User.Email); err != nil {
		return nil, err
	}

	if err := s.actionService.Create(ctx, &models.UserAction{
		UserID: userID,
		Action: models.ActionUserVerification,
	}); err != nil {
		return nil, err
	}

	session, err := s.sessionService.Create(ctx, userID, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return toSessionResponse(session), nil
}

func (s *AuthService) SignIn(ctx context.Context, dto *dtos.SignInDto) (*responses.SessionResponse, error) {
	user, err := s.userService.GetByUsername(ctx, dto.Username)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if !util.ValidatePassword(dto.Password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	session, err := s.sessionService.Create(ctx, user.ID, &dto.SessionData)
	if err != nil {
		return nil, err
	}

	return toSessionResponse(session), nil
}

func (s *AuthService) SignOut(ctx context.Context, sessionID *uuid.UUID) error {
	return s.sessionService.Close(ctx, sessionID)
}

func (s *AuthService) Refresh(ctx context.Context, session *models.Session, dto *dtos.SessionDataDto) (*responses.SessionResponse, error) {
	if session.Device != dto.Device || session.UserAgent != dto.UserAgent {
		if err := s.sessionService.Close(ctx, &session.ID); err != nil {
			return nil, err
		}
		return nil, models.ErrSessionInvalid
	}
	updated, err := s.sessionService.Update(ctx, session)
	if err != nil {
		return nil, err
	}
	return toSessionResponse(updated), nil
}

func (s *AuthService) ForgotPasswordRequest(ctx context.Context, email string) error {
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	if err := s.verificationService.GenerateAndSend(ctx, user.ID, user.Email); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ForgotPasswordConfirm(
	ctx context.Context,
	token string,
	dto *dtos.SessionDataDto,
) (*responses.SessionResponse, error) {
	verification, err := s.verificationService.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionService.Create(ctx, verification.UserID, dto)
	if err != nil {
		return nil, err
	}

	return toSessionResponse(session), nil
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

	if err := s.sessionService.CloseAll(ctx, session.User.ID); err != nil {
		return err
	}

	return nil
}

func toSessionResponse(s *models.Session) *responses.SessionResponse {
	return &responses.SessionResponse{
		SessionId:    s.ID.String(),
		UserId:       s.User.ID.String(),
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		Status:       s.Status.String(),
		Device:       s.Device,
		Ip:           s.IpAddress,
		Agent:        s.UserAgent,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
