package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type AuthHandler struct {
	aService *service.AuthService
	sService *service.SessionService
}

func NewAuthHandler(
	aService *service.AuthService,
	sService *service.SessionService,
) *AuthHandler {
	return &AuthHandler{
		aService: aService,
		sService: sService,
	}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/signIn", h.signIn)
		auth.Post("/signUp", h.signUp)

		auth.Group(func(protected chi.Router) {
			protected.Use(middlewares.RequiresSession(h.sService.GetByID))
			protected.Post("/signOut", h.signOut)
			protected.Post("/ping", h.ping)
			protected.Post("/refresh", h.refresh)
		})

		auth.Post("/forgot/password/confirm", h.forgotPasswordConfirm)
	})
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignInDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	session, err := h.aService.SignIn(r.Context(), &dto)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(session.Response(), "sign in successful"))
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignUpDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	session, err := h.aService.SignUp(r.Context(), &dto)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(session.Response(), "sign up successful"))
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

	if err := h.aService.SignOut(r.Context(), session.ID); err != nil {
		responses.HandleReqError(w, err)
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

	res, err := h.aService.Refresh(r.Context(), session, &dto)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(res, "Session refreshed"))
}

func (h *AuthHandler) forgotPasswordConfirm(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ForgotPasswordConfirmDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	parsedID, err := uuid.Parse(dto.ActionID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	session, err := h.aService.ForgotPasswordConfirm(r.Context(), parsedID, &dto.SessionData)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(session, "Success"))
}
