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
	userService *service.UserService
	pHandler    *UserPreferencesHandler
}

func NewUserMeHandler(userService *service.UserService, pHandler *UserPreferencesHandler) *UserMeHandler {
	return &UserMeHandler{
		userService: userService,
		pHandler:    pHandler,
	}
}

func (h *UserMeHandler) Routes(r chi.Router) {
	r.Route("/me", func(me chi.Router) {
		me.Put("/", h.updateMe)
		me.Delete("/", h.deleteMe)
		// TODO:  email update routes

		h.pHandler.Routes(me)
	})
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
