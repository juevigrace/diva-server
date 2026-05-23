package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserPreferencesHandler struct {
	service        *service.UserPreferencesService
	sessionService *service.SessionService
}

func NewUserPreferencesHandler(svc *service.UserPreferencesService, sessionService *service.SessionService) *UserPreferencesHandler {
	return &UserPreferencesHandler{service: svc, sessionService: sessionService}
}

func (h *UserPreferencesHandler) Routes(r chi.Router) {
	r.Route("/preferences", func(pref chi.Router) {
		pref.Get("/", h.getByUser)
		pref.Get("/{id}", h.getByID)
		pref.Get("/by-device/{device}", h.getByDevice)
		pref.Post("/", h.createPreferences)
		pref.Put("/", h.updatePreferences)

		pref.Group(func(admin chi.Router) {
			admin.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
			admin.Delete("/{id}", h.deletePreferences)
		})
	})
}

func (h *UserPreferencesHandler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return false
	}
	if session.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "admin access required"))
		return false
	}
	return true
}

func (h *UserPreferencesHandler) getByUser(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	prefs, err := h.service.GetByUser(r.Context(), session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res := make([]*responses.UserPreferencesResponse, len(prefs))
	for i, p := range prefs {
		res[i] = p.Response(&session.User.ID)
	}

	responses.WriteJSON(w, responses.RespondOk(res, "Preferences retrieved"))
}

func (h *UserPreferencesHandler) getByID(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	pref, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(pref.Response(&session.User.ID), "Preference retrieved"))
}

func (h *UserPreferencesHandler) getByDevice(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	device := chi.URLParam(r, "device")
	if device == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "device is required"))
		return
	}

	pref, err := h.service.GetByDevice(r.Context(), device)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(pref.Response(&session.User.ID), "Preference retrieved"))
}

func (h *UserPreferencesHandler) createPreferences(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.CreateUserPreferencesDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	dto.Device = session.Device

	if err := h.service.Create(r.Context(), session.User.ID, &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "Preferences created"))
}

func (h *UserPreferencesHandler) updatePreferences(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UpdateUserPreferencesDto
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

func (h *UserPreferencesHandler) deletePreferences(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Preference deleted"))
}
