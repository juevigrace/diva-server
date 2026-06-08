package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/errs"
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

func (h *UserPreferencesHandler) UserRoutes(r chi.Router) {
	r.Route("/preferences", func(pref chi.Router) {
		pref.With(middlewares.RequireResourceOwner(
			"uid",
			func(_ context.Context, reqid, resid uuid.UUID) (any, bool) {
				return nil, reqid == resid
			},
			models.PERMISSION_USER_PERMISSIONS_READ,
		)).Get("/", h.getByUser)

		pref.With(
			middlewares.RequirePermission(models.PERMISSION_USERS_PREFERENCES_WRITE),
			middlewares.RequireResourceOwner(
				"uid",
				func(_ context.Context, reqid, resid uuid.UUID) (any, bool) {
					return nil, reqid == resid
				},
				models.PERMISSION_USERS_PREFERENCES_WRITE,
			),
		).Post("/", h.createPreferences)
	})
}

func (h *UserPreferencesHandler) Routes(r chi.Router) {
	r.Route("/preferences", func(p chi.Router) {
		p.Route("/{pid}", func(pid chi.Router) {
			pid.With(middlewares.RequireResourceOwner(
				"pid",
				func(ctx context.Context, reqid, resid uuid.UUID) (any, bool) {
					pref, err := h.upService.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if pref.UserID != reqid {
						return nil, false
					}
					return pref, true
				},
				models.PERMISSION_USERS_PREFERENCES_READ,
			)).Get("/", h.getByID)

			pid.With(middlewares.RequireResourceOwner(
				"pid",
				func(ctx context.Context, reqid, resid uuid.UUID) (any, bool) {
					pref, err := h.upService.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if pref.UserID != reqid {
						return nil, false
					}
					return pref, true
				},
				models.PERMISSION_USERS_PREFERENCES_WRITE,
			)).Put("/", h.updatePreferences)
		})
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
		res[i] = p.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(res, "preferences retrieved"))
}

func (h *UserPreferencesHandler) getByID(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	pref, ok := rc.Cache["pid"].(*models.UserPreferences)
	if !ok {
		pid, err := middlewares.GetUUIDFromURL(r, "pid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
		pref, err = h.upService.GetByID(r.Context(), pid)
		if err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondOk(pref.Response(), "preference retrieved"))
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

	if err = h.upService.Create(r.Context(), uid, &dto); err != nil {
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
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	var pid uuid.UUID
	prefs, ok := rc.Cache["pid"].(*models.UserPreferences)
	if !ok {
		pid, err = middlewares.GetUUIDFromURL(r, "pid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	} else {
		pid = prefs.ID
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
