package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserPreferencesHandler struct {
	service *service.UserPreferencesService
}

func NewUserPreferencesHandler(svc *service.UserPreferencesService) *UserPreferencesHandler {
	return &UserPreferencesHandler{service: svc}
}

func (h *UserPreferencesHandler) Routes(r chi.Router) {
	r.Route("/preferences", func(pref chi.Router) {
		// TODO: get handler
		pref.Post("/", h.createPreferences)
		pref.Put("/", h.updatePreferences)
	})
}

func (h *UserPreferencesHandler) createPreferences(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.UserPreferencesDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Create(r.Context(), session.User.ID, &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "Preferences created"))
}

func (h *UserPreferencesHandler) updatePreferences(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserPreferencesDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Update(r.Context(), &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Preferences updated"))
}
