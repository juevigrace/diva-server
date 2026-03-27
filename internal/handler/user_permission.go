package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
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
	r.Route("/permissions", func(perms chi.Router) {
		perms.Get("/", h.getPermission)
		// TODO: add perm checks
		perms.Post("/", h.createPermission)
		perms.Put("/", h.updatePermission)
		perms.Delete("/", h.deletePermission)
	})
}

func (h *UserPermissionHandler) getPermission(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	dto, err := h.service.Get(r.Context(), session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(dto, "Permission retrieved"))
}

func (h *UserPermissionHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	var dto dtos.UserPermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Create(r.Context(), &dto); err != nil {
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
