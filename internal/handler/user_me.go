package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserMeHandler struct {
	actionService       *service.UserActionsService
	userService         *service.UserService
	verificationService *service.VerificationService
	pHandler            *UserPreferencesHandler
	aHandler            *UserActionsHandler
}

func NewUserMeHandler(
	actionService *service.UserActionsService,
	userService *service.UserService,
	verificationService *service.VerificationService,
	pHandler *UserPreferencesHandler,
	aHandler *UserActionsHandler,
) *UserMeHandler {
	return &UserMeHandler{
		actionService:       actionService,
		userService:         userService,
		verificationService: verificationService,
		pHandler:            pHandler,
		aHandler:            aHandler,
	}
}

func (h *UserMeHandler) Routes(r chi.Router) {
	r.Route("/me", func(me chi.Router) {
		me.Put("/", h.updateMe)
		me.Delete("/", h.deleteMe)
		me.Route("/email", func(verify chi.Router) {
			verify.Post("/verify", h.verifyEmail)

			// TODO:  email update routes
		})

		h.pHandler.Routes(me)
		h.aHandler.Routes(me)
	})
}

func (h *UserMeHandler) verifyEmail(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.EmailTokenDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if _, err := h.verificationService.Verify(r.Context(), dto.Token); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.userService.UpdateVerified(r.Context(), &session.User.ID); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.actionService.Delete(r.Context(), &models.UserAction{
		UserID: session.User.ID,
		Action: models.ActionUserVerification,
	}); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "user verified"))
}

func (h *UserMeHandler) updateMe(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.UpdateProfileDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.userService.UpdateProfile(r.Context(), session.User.ID, &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Profile updated"))
}

func (h *UserMeHandler) deleteMe(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	if err := h.userService.Delete(r.Context(), session.User.ID); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "User deleted"))
}
