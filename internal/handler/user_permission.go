package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserPermissionHandler struct {
	service *service.UserPermissionService
}

func NewUserPermissionHandler(svc *service.UserPermissionService) *UserPermissionHandler {
	return &UserPermissionHandler{service: svc}
}

// TODO: implement better internal logic to a better management of permissions
func (h *UserPermissionHandler) UserRoutes(r chi.Router) {
	r.Route("/permissions", func(perms chi.Router) {
		perms.Group(func(rg chi.Router) {
			rg.Use(middlewares.RequireResourceOwner(
				"uid",
				func(_ context.Context, reqid, resid uuid.UUID) (any, bool) {
					return nil, reqid == resid
				},
				models.PERMISSION_USER_PERMISSIONS_READ,
			))
			rg.Get("/", h.getByUser)
			rg.Get("/{pid}", h.getOneByUser)
		})

		perms.Group(func(admin chi.Router) {
			admin.Use(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
				middlewares.RequirePermission(models.PERMISSION_USER_PERMISSIONS_WRITE),
			)
			admin.Delete("/{pid}", h.deletePermission)
			admin.Delete("/", h.deleteByUser)
		})
	})

}

func (h *UserPermissionHandler) Routes(r chi.Router) {
	r.Route("/permissions", func(perms chi.Router) {
		perms.Use(
			middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			middlewares.RequirePermission(models.PERMISSION_USER_PERMISSIONS_WRITE),
		)
		perms.Post("/", h.createPermission)
		perms.Put("/", h.updatePermission)
	})
}

func (h *UserPermissionHandler) getByUser(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	perms, err := h.service.GetByUser(r.Context(), uid)
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

	perm, err := h.service.GetOneByUser(r.Context(), uid, pid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(), "permission retrieved"))
}

func (h *UserPermissionHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	permissionID, err := uuid.Parse(dto.PermissionId)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	userID, err := uuid.Parse(dto.UserId)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
		return
	}

	perm := &models.UserPermission{
		Permission: models.Permission{ID: permissionID},
		GrantedBy:  &session.User.ID,
		Granted:    dto.Granted,
		GrantedAt:  time.Now().UTC().UnixMilli(),
		ExpiresAt:  dto.ExpiresAt,
	}

	if err := h.service.Create(r.Context(), userID, perm); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "permission created"))
}

func (h *UserPermissionHandler) updatePermission(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Update(r.Context(), &dto); err != nil {
		responses.HandleReqError(w, err)
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

	if err := h.service.Delete(r.Context(), uid, pid); err != nil {
		responses.HandleReqError(w, err)
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

	if err := h.service.DeleteByUser(r.Context(), uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "permissions deleted"))
}
