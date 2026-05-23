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

type UserPermissionHandler struct {
	service        *service.UserPermissionService
	sessionService *service.SessionService
}

func NewUserPermissionHandler(svc *service.UserPermissionService, sessionService *service.SessionService) *UserPermissionHandler {
	return &UserPermissionHandler{service: svc, sessionService: sessionService}
}

func (h *UserPermissionHandler) Routes(r chi.Router) {
	r.Route("/permissions", func(perms chi.Router) {
		perms.Get("/", h.getPermission)
		perms.Post("/", h.createPermission)
		perms.Put("/", h.updatePermission)
		perms.Delete("/", h.deletePermission)

		perms.Group(func(admin chi.Router) {
			admin.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
			admin.Get("/{permID}", h.getOneByUser)
			admin.Delete("/user/{userId}", h.deleteByUser)
		})
	})
}

func (h *UserPermissionHandler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
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

func (h *UserPermissionHandler) getPermission(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	dto, err := h.service.GetByUser(r.Context(), session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(dto, "Permission retrieved"))
}

func (h *UserPermissionHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.UserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Create(r.Context(), session, &dto); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "Permission created"))
}

func (h *UserPermissionHandler) updatePermission(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Update(r.Context(), &dto); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Permission updated"))
}

func (h *UserPermissionHandler) deletePermission(w http.ResponseWriter, r *http.Request) {
	var dto dtos.DeleteUserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Delete(r.Context(), &dto); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Permission deleted"))
}

func (h *UserPermissionHandler) getOneByUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	userIDParam := chi.URLParam(r, "userId")
	permIDParam := chi.URLParam(r, "permID")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid user id"))
		return
	}

	permID, err := uuid.Parse(permIDParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid permission id"))
		return
	}

	up, err := h.service.GetOneByUser(r.Context(), userID, permID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(up.Response(&userID), "Permission retrieved"))
}

func (h *UserPermissionHandler) deleteByUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	userIDParam := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid user id"))
		return
	}

	if err := h.service.DeleteByUser(r.Context(), userID); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Permissions deleted"))
}
