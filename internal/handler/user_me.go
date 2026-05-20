package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserMeHandler struct {
	uService  *service.UserService
	upService *service.UserProfileService
	pHandler  *UserPreferencesHandler
	aHandler  *UserActionsHandler
}

func NewUserMeHandler(
	uService *service.UserService,
	upService *service.UserProfileService,
	pHandler *UserPreferencesHandler,
	aHandler *UserActionsHandler,
) *UserMeHandler {
	return &UserMeHandler{
		uService:  uService,
		upService: upService,
		pHandler:  pHandler,
		aHandler:  aHandler,
	}
}

func (h *UserMeHandler) Routes(r chi.Router) {
	r.Route("/me", func(me chi.Router) {
		me.Get("/", h.getMe)
		me.Put("/", h.updateMe)
		me.Delete("/", h.deleteMe)
		me.Patch("/email", func(w http.ResponseWriter, r *http.Request) {})

		h.pHandler.Routes(me)
		h.aHandler.Routes(me)
	})
}

func (h *UserMeHandler) getMe(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	user, err := h.uService.GetByID(r.Context(), session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(user.Response(), "Good"))
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

	if err := h.upService.Update(r.Context(), session.User.ID, &dto); err != nil {
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

	if err := h.uService.Delete(r.Context(), session.User.ID); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "User deleted"))
}
