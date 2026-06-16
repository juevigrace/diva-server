package preferences

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/core/user/permissions"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type UserPreferencesHandler struct {
	upRepo  *permissions.UserPermissionRepo
	uprRepo *UserPreferencesRepo
}

func NewUserPreferencesHandler(
	upRepo *permissions.UserPermissionRepo,
	uprRepo *UserPreferencesRepo,
) *UserPreferencesHandler {
	return &UserPreferencesHandler{
		upRepo:  upRepo,
		uprRepo: uprRepo,
	}
}

func (h *UserPreferencesHandler) UserRoutes(r chi.Router) {
	r.Route("/preferences", func(pref chi.Router) {
		pref.With(middlewares.RequireResourceOwner(
			&middlewares.RequireOwnerParams{
				UrlParams: []string{"uid"},
				Perms:     []models.PermissionAction{models.PERMISSION_USER_PERMISSIONS_READ},
			},
			func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
				resid, err := uuid.Parse(resParams[0])
				if err != nil {
					return nil, false
				}
				if reqid != resid {
					return nil, false
				}
				return map[string]any{"uid": resid}, true
			},
		)).Get("/", h.getByUser)

		pref.With(
			middlewares.RequirePermission(models.PERMISSION_USERS_PREFERENCES_WRITE),
			middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"uid"},
					Perms:     []models.PermissionAction{models.PERMISSION_USERS_PREFERENCES_WRITE},
				},
				func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					if reqid != resid {
						return nil, false
					}
					return nil, true
				},
			),
		).Post("/", h.createPreferences)
	})
}

func (h *UserPreferencesHandler) Routes(r chi.Router) {
	r.Route("/preferences", func(p chi.Router) {
		p.Route("/{pid}", func(pid chi.Router) {
			pid.With(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"pid"},
					Perms:     []models.PermissionAction{models.PERMISSION_USERS_PREFERENCES_READ},
				},
				func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					pref, err := h.uprRepo.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if pref.UserID != reqid {
						return nil, false
					}
					return map[string]any{"pid": pref}, true
				},
			)).Get("/", h.getByID)

			pid.With(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"pid"},
					Perms:     []models.PermissionAction{models.PERMISSION_USERS_PREFERENCES_WRITE},
				},
				func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					pref, err := h.uprRepo.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if pref.UserID != reqid {
						return nil, false
					}
					return map[string]any{"pid": pref}, true
				},
			)).Put("/", h.updatePreferences)
		})
	})
}

func (h *UserPreferencesHandler) getByUser(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	uid, ok := rc.Cache["uid"].(uuid.UUID)
	if !ok {
		uid, err = middlewares.GetUUIDFromURL(r, "uid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	}

	prefs, err := h.uprRepo.GetByUser(r.Context(), uid)
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
		pref, err = h.uprRepo.GetByID(r.Context(), pid)
		if err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondOk(pref.Response(), "preference retrieved"))
}

func (h *UserPreferencesHandler) createPreferences(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	uid, ok := rc.Cache["uid"].(uuid.UUID)
	if !ok {
		uid, err = middlewares.GetUUIDFromURL(r, "uid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
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

	if err = h.uprRepo.Create(r.Context(), rc.Session, uid, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
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

	if err = h.uprRepo.Update(r.Context(), pid, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "preferences updated"))
}
