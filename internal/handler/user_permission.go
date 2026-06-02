package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserPermissionHandler struct {
	service *service.UserPermissionService
}

func NewUserPermissionHandler(svc *service.UserPermissionService) *UserPermissionHandler {
	return &UserPermissionHandler{service: svc}
}

func (h *UserPermissionHandler) Routes(r chi.Router) {
	r.Route("/{uid}/permissions", func(uid chi.Router) {
		uid.Get("/", h.getAll)
		uid.Route("/{pid}", func(pid chi.Router) {
			pid.Get("/", h.getOneByUser)
			pid.Delete("/", h.deletePermission)
		})
		uid.Delete("/", h.deleteByUser)
	})
	r.Route("/permissions", func(perms chi.Router) {
		perms.Post("/", h.createPermission)
		perms.Put("/", h.updatePermission)
	})
}

func (h *UserPermissionHandler) getAll(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	perms, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*[]models.UserPermission, error) {
			dbPerms, err := h.service.GetByUser(r.Context(), uid)
			if err != nil {
				return nil, err
			}

			return &dbPerms, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	res := make([]*responses.UserPermissionResponse, len(*perms))
	for i, p := range *perms {
		res[i] = p.Response(&uid)
	}

	responses.WriteJSON(w, responses.RespondOk(res, "permission retrieved"))
}

func (h *UserPermissionHandler) getOneByUser(w http.ResponseWriter, r *http.Request) {
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

	perm, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*models.UserPermission, error) {
			return h.service.GetOneByUser(r.Context(), uid, pid)
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(&uid), "permission retrieved"))
}

func (h *UserPermissionHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	_, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.UserPermission, error) {
			var dto dtos.UserPermissionDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			permissionID, err := uuid.Parse(dto.PermissionId)
			if err != nil {
				return nil, err
			}

			userID, err := uuid.Parse(dto.UserId)
			if err != nil {
				return nil, err
			}

			grantedAt := new(int64)
			if dto.Granted {
				*grantedAt = time.Now().UTC().UnixMilli()
			}

			perm := &models.UserPermission{
				Permission: models.Permission{ID: permissionID},
				GrantedBy:  &session.User.ID,
				Granted:    dto.Granted,
				GrantedAt:  grantedAt,
				ExpiresAt:  dto.ExpiresAt,
			}

			if err := h.service.Create(r.Context(), userID, perm); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "permission created"))
}

func (h *UserPermissionHandler) updatePermission(w http.ResponseWriter, r *http.Request) {
	_, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.UserPermission, error) {
			var dto dtos.UserPermissionDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.service.Update(r.Context(), &dto); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "permission updated"))
}

func (h *UserPermissionHandler) deletePermission(w http.ResponseWriter, r *http.Request) {
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.UserPermission, error) {
			if err := h.service.Delete(r.Context(), uid, pid); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "permission deleted"))
}

func (h *UserPermissionHandler) deleteByUser(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.UserPermission, error) {
			if err := h.service.DeleteByUser(r.Context(), uid); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "permissions deleted"))
}
