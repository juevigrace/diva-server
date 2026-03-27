package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type AuthHandler struct {
	authService    *service.AuthService
	sessionService *service.SessionService
}

func NewAuthHandler(authService *service.AuthService, sessionService *service.SessionService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionService: sessionService,
	}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/signIn", h.signIn)
		auth.Post("/signUp", h.signUp)

		auth.Group(func(protected chi.Router) {
			protected.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
			protected.Post("/signOut", h.signOut)
			protected.Post("/ping", h.ping)
			protected.Post("/refresh", h.refresh)
		})

		auth.Route("/forgot", func(forgot chi.Router) {
			forgot.Route("/password", func(pass chi.Router) {
				pass.Group(func(upPass chi.Router) {
					upPass.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
					upPass.Patch("/", h.forgotPasswordUpdate)
				})
			})
		})
	})
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignInDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	session, err := h.authService.SignIn(r.Context(), &dto)
	if err != nil {
		res := new(responses.APIResponse)
		if errors.Is(models.ErrInvalidCredentials, err) {
			res = responses.RespondUnauthorized(nil, err.Error())
		} else {
			res = responses.RespondBadRequest(nil, err.Error())
		}
		responses.WriteJSON(w, res)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(session, "Sign in successful"))
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignUpDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	session, err := h.authService.SignUp(r.Context(), &dto)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(session, "Sign up successful"))
}

func (h *AuthHandler) signOut(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.SessionDataDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.authService.SignOut(r.Context(), &session.ID); err != nil {
		if errors.Is(err, models.ErrSessionInvalid) {
			responses.WriteJSON(w, responses.RespondUnauthorized(nil, err.Error()))
			return
		}
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
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
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res, err := h.authService.Refresh(r.Context(), session, &dto)
	if err != nil {
		if errors.Is(err, models.ErrSessionInvalid) {
			responses.WriteJSON(w, responses.RespondUnauthorized(nil, err.Error()))
			return
		}
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(res, "Session refreshed"))
}

func (h *AuthHandler) forgotPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.UpdatePasswordDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.authService.ForgotPasswordUpdate(r.Context(), session, &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Password updated successfully"))
}
