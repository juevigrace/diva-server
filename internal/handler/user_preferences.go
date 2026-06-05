package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserPreferencesHandler struct {
	upService  *service.UserPreferencesService
	uprService *service.UserPermissionService
}

func NewUserPreferencesHandler(
	upService *service.UserPreferencesService,
	uprService *service.UserPermissionService,
) *UserPreferencesHandler {
	return &UserPreferencesHandler{
		upService:  upService,
		uprService: uprService,
	}
}

func (h *UserPreferencesHandler) Routes(r chi.Router) {
	r.Route("/preferences", func(pref chi.Router) {
		pref.Get("/", h.getByUser)
		pref.Route("/{pid}", func(uid chi.Router) {
			uid.Get("/", h.getByID)
			uid.Put("/", h.updatePreferences)
		})
		pref.Post("/", h.createPreferences)
	})
}

func (h *UserPreferencesHandler) getByUser(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	prefs, err := h.upService.GetByUser(r.Context(), uid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	res := make([]*responses.UserPreferencesResponse, len(prefs))
	for i, p := range prefs {
		res[i] = p.Response(&uid)
	}

	responses.WriteJSON(w, responses.RespondOk(res, "preferences retrieved"))
}

func (h *UserPreferencesHandler) getByID(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	pref, err := h.upService.GetByID(r.Context(), pid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(pref.Response(&uid), "preference retrieved"))
}

func (h *UserPreferencesHandler) createPreferences(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
		return
	}

	var dto dtos.CreateUserPreferencesDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		return
	}

	dto.Device = session.Device

	if err = h.upService.Create(r.Context(), session.User.ID, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.ID == uid {
		if err := h.uprService.Delete(r.Context(), session.User.ID, session.User.Permissions[models.PERMISSION_USERS_PREFERENCES_WRITE].Permission.ID); err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "preferences created"))
}

func (h *UserPreferencesHandler) updatePreferences(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	var dto dtos.UpdateUserPreferencesDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		return
	}

	if err = h.upService.Update(r.Context(), pid, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "preferences updated"))
}
