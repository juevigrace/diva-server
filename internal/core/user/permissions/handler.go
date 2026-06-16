package permissions

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type UserPermissionHandler struct {
	upRepo *UserPermissionRepo
}

func NewUserPermissionHandler(upRepo *UserPermissionRepo) *UserPermissionHandler {
	return &UserPermissionHandler{upRepo: upRepo}
}

func (h *UserPermissionHandler) UserRoutes(r chi.Router) {
	r.Route("/permissions", func(perms chi.Router) {
		perms.With(middlewares.RequireResourceOwner(
			&middlewares.RequireOwnerParams{
				UrlParams: []string{"uid"},
				Perms:     []models.PermissionAction{models.PERMISSION_USER_PERMISSIONS_READ},
			},
			func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
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

		perms.Route("/{pid}", func(pid chi.Router) {
			pid.With(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"uid", "pid"},
					Perms:     []models.PermissionAction{models.PERMISSION_USER_PERMISSIONS_READ},
				},
				func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					uid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					if reqid != uid {
						return nil, false
					}
					pid, err := uuid.Parse(resParams[1])
					if err != nil {
						return nil, false
					}
					return map[string]any{"uid": uid, "pid": pid}, true
				},
			)).Get("/", h.getOneByUser)

			pid.Group(func(admin chi.Router) {
				admin.Use(
					middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
					middlewares.RequireResourceOwner(
						&middlewares.RequireOwnerParams{
							UrlParams: []string{"uid", "pid"},
						},
						func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
							uid, err := uuid.Parse(resParams[0])
							if err != nil {
								return nil, false
							}
							pid, err := uuid.Parse(resParams[1])
							if err != nil {
								return nil, false
							}
							dbPerm, err := h.upRepo.GetOneByPermID(ctx, uid, pid)
							if err != nil {
								return nil, false
							}
							if dbPerm.GrantedBy != nil && reqid != *dbPerm.GrantedBy {
								return nil, false
							}
							return map[string]any{"pid": pid, "uid": uid}, true
						},
					),
					middlewares.RequirePermission(models.PERMISSION_USER_PERMISSIONS_WRITE),
				)
				admin.Put("/", h.updatePermission)
				admin.Delete("/", h.deletePermission)
			})
		})

		perms.With(
			middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			middlewares.RequirePermission(models.PERMISSION_USER_PERMISSIONS_WRITE),
		).Post("/", h.createPermission)
	})
}

func (h *UserPermissionHandler) getByUser(w http.ResponseWriter, r *http.Request) {
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

	perms, err := h.upRepo.GetByUser(r.Context(), uid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	res := make([]*responses.UserPermissionResponse, len(perms))
	for i, p := range perms {
		res[i] = p.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(res, "permission retrieved"))
}

func (h *UserPermissionHandler) getOneByUser(w http.ResponseWriter, r *http.Request) {
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

	pid, ok := rc.Cache["pid"].(uuid.UUID)
	if !ok {
		pid, err = middlewares.GetUUIDFromURL(r, "pid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	}

	perm, err := h.upRepo.GetOneByPermID(r.Context(), uid, pid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(), "permission retrieved"))
}

func (h *UserPermissionHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if rc.Session.User.Role != models.ROLE_ADMIN && rc.Session.User.ID == uid {
		responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
		return
	}

	var dto dtos.CreateUserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	permission := models.PermissionActionFromString(dto.PermissionAction)
	if permission == models.PERMISSION_NONE {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, errs.ErrPermissionNotFound.Error()))
		return
	}

	if err := h.upRepo.CreateByName(r.Context(), permission, &rc.Session.User, dto.Granted, dto.ExpiresAt, uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "permission created"))
}

func (h *UserPermissionHandler) updatePermission(w http.ResponseWriter, r *http.Request) {
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

	if rc.Session.User.Role != models.ROLE_ADMIN && rc.Session.User.ID == uid {
		responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
		return
	}

	pid, ok := rc.Cache["pid"].(uuid.UUID)
	if !ok {
		pid, err = middlewares.GetUUIDFromURL(r, "pid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	}

	var dto dtos.UpdateUserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.upRepo.Update(r.Context(), uid, pid, dto.Granted, dto.ExpiresAt); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "permission updated"))
}

func (h *UserPermissionHandler) deletePermission(w http.ResponseWriter, r *http.Request) {
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

	if rc.Session.User.Role != models.ROLE_ADMIN && rc.Session.User.ID == uid {
		responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
		return
	}

	pid, ok := rc.Cache["pid"].(uuid.UUID)
	if !ok {
		pid, err = middlewares.GetUUIDFromURL(r, "pid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	}

	if err := h.upRepo.Delete(r.Context(), uid, pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "permission deleted"))
}
