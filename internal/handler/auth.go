package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/repo"
	"github.com/juevigrace/diva-server/internal/util"
)

const (
	ActionEmailVerification = "EMAIL_VERIFICATION"
)

type AuthHandler struct {
	uRepo        *repo.UserRepository
	sRepo        *repo.SessionRepository
	verification *repo.VerificationRepository
	mail         *mail.Client
}

func NewAuthHandler(
	uRepo *repo.UserRepository,
	sRepo *repo.SessionRepository,
	verification *repo.VerificationRepository,
	mail *mail.Client,
) *AuthHandler {
	return &AuthHandler{
		uRepo:        uRepo,
		sRepo:        sRepo,
		verification: verification,
		mail:         mail,
	}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/signIn", h.signIn)
		auth.Post("/signUp", h.signUp)

		auth.Group(func(protected chi.Router) {
			protected.Use(middlewares.SessionMiddleware(h.sRepo.GetByID))
			protected.Post("/signOut", h.signOut)
			protected.Post("/ping", h.ping)
			protected.Post("/refresh", h.refresh)

			protected.Route("/forgot/password", func(forgot chi.Router) {
				forgot.Patch("/", h.forgotPasswordUpdate)
			})
		})

		auth.Route("/forgot/password", func(forgot chi.Router) {
			forgot.Post("/request", h.forgotPasswordRequest)
			forgot.Post("/confirm", h.forgotPasswordConfirm)
		})
	})
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignInDto
	if _, err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	user, err := h.uRepo.GetByUsername(r.Context(), dto.Username)
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "invalid credentials"))
		return
	}

	if !user.UserVerified {
		verification, verErr := h.verification.Create(r.Context(), user.ID)
		if verErr == nil {
			go func() {
				_ = h.mail.SendVerificationEmail(r.Context(), user.Email, verification)
			}()
		}
	}

	if !util.ValidatePassword(dto.Password, *user.PasswordHash) {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "invalid credentials"))
		return
	}

	sessionID := uuid.New()
	sessionDto := &dtos.SessionDataDto{
		Device:    dto.SessionData.Device,
		UserAgent: dto.SessionData.UserAgent,
		IpAddress: dto.SessionData.IpAddress,
	}

	session := &models.Session{
		ID:        sessionID,
		User:      user.ID,
		Device:    sessionDto.Device,
		UserAgent: sessionDto.UserAgent,
		IpAddress: sessionDto.IpAddress,
	}

	if err := h.sRepo.Create(r.Context(), sessionID, user.ID, sessionDto); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to create session"))
		return
	}

	actions := []*responses.ActionResponse{}
	if !user.UserVerified {
		actions = append(actions, &responses.ActionResponse{ActionName: ActionEmailVerification})
	}

	responses.WriteJSON(w, responses.RespondOk(&responses.AuthResponse{
		Session: h.toSessionResponse(session),
		Actions: actions,
	}, "Sign in successful"))
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignUpDto
	if _, err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err := h.uRepo.GetByUsername(r.Context(), dto.User.Username)
	if err == nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "username already exists"))
		return
	}

	userID, err := h.uRepo.Create(r.Context(), dto.User)
	if err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to create user"))
		return
	}

	go func() {
		verification, verErr := h.verification.Create(r.Context(), userID)
		if verErr != nil {
			return
		}
		_ = h.mail.SendVerificationEmail(r.Context(), dto.User.Email, verification)
	}()

	sessionID := uuid.New()
	sessionDto := &dtos.SessionDataDto{
		Device:    dto.SessionData.Device,
		UserAgent: dto.SessionData.UserAgent,
		IpAddress: dto.SessionData.IpAddress,
	}

	session := &models.Session{
		ID:        sessionID,
		User:      userID,
		Device:    sessionDto.Device,
		UserAgent: sessionDto.UserAgent,
		IpAddress: sessionDto.IpAddress,
	}

	if err := h.sRepo.Create(r.Context(), sessionID, userID, sessionDto); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to create session"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(&responses.AuthResponse{
		Session: h.toSessionResponse(session),
		Actions: []*responses.ActionResponse{
			{ActionName: ActionEmailVerification},
		},
	}, "Sign up successful"))
}

func (h *AuthHandler) signOut(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	session.Status = models.SESSION_CLOSED
	if err := h.sRepo.Update(r.Context(), session.ID, session.User, &dtos.SessionDataDto{
		Device:    session.Device,
		UserAgent: session.UserAgent,
		IpAddress: session.IpAddress,
	}); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to sign out"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Sign out successful"))
}

func (h *AuthHandler) ping(w http.ResponseWriter, r *http.Request) {
	responses.WriteJSON(w, responses.RespondOk(nil, "Pong"))
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.SessionDataDto
	if _, err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if session.Device != dto.Device || session.UserAgent != dto.UserAgent {
		_ = h.sRepo.Delete(r.Context(), session.ID)
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "invalid session"))
		return
	}

	if err := h.sRepo.Update(r.Context(), session.ID, session.User, &dto); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to refresh session"))
		return
	}

	updatedSession, err := h.sRepo.GetByID(r.Context(), session.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to get session"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(h.toSessionResponse(updatedSession), "Session refreshed"))
}

func (h *AuthHandler) forgotPasswordRequest(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ForgotPasswordRequestDto
	if _, err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	user, err := h.uRepo.GetByEmail(r.Context(), dto.Email)
	if err != nil {
		responses.WriteJSON(w, responses.RespondOk(nil, "If the email exists, a verification code has been sent"))
		return
	}

	verification, verErr := h.verification.Create(r.Context(), user.ID)
	if verErr != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to create verification"))
		return
	}

	go func() {
		_ = h.mail.SendVerificationEmail(r.Context(), user.Email, verification)
	}()

	responses.WriteJSON(w, responses.RespondOk(nil, "If the email exists, a verification code has been sent"))
}

func (h *AuthHandler) forgotPasswordConfirm(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ForgotPasswordConfirmDto
	if _, err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	verification, err := h.verification.Verify(r.Context(), dto.Token)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid or expired token"))
		return
	}

	sessionID := uuid.New()
	sessionDto := &dtos.SessionDataDto{
		Device:    dto.SessionData.Device,
		UserAgent: dto.SessionData.UserAgent,
		IpAddress: dto.SessionData.IpAddress,
	}

	session := &models.Session{
		ID:        sessionID,
		User:      verification.UserID,
		Device:    sessionDto.Device,
		UserAgent: sessionDto.UserAgent,
		IpAddress: sessionDto.IpAddress,
	}

	if err := h.sRepo.Create(r.Context(), sessionID, verification.UserID, sessionDto); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to create session"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(h.toSessionResponse(session), "Password reset session created"))
}

func (h *AuthHandler) forgotPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.UpdatePasswordDto
	if _, err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.uRepo.UpdatePassword(r.Context(), session.User, &dto); err != nil {
		if err.Error() == "new password cannot be the same as current password" {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "failed to update password"))
		return
	}

	session.Status = models.SESSION_CLOSED
	_ = h.sRepo.Delete(r.Context(), session.ID)

	responses.WriteJSON(w, responses.RespondOk(nil, "Password updated successfully"))
}

func (h *AuthHandler) toSessionResponse(session *models.Session) *responses.SessionResponse {
	return &responses.SessionResponse{
		SessionId:    session.ID.String(),
		UserId:       session.User.String(),
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		Status:       session.Status.String(),
		Device:       session.Device,
		Ip:           session.IpAddress,
		Agent:        session.UserAgent,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}
}
