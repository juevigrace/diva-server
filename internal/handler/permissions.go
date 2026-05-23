package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type PermissionsHandler struct {
	service        *service.PermissionService
	sessionService *service.SessionService
}

func NewPermissionsHandler(svc *service.PermissionService, sessionService *service.SessionService) *PermissionsHandler {
	return &PermissionsHandler{
		service:        svc,
		sessionService: sessionService,
	}
}

func (h *PermissionsHandler) Routes(r chi.Router) {
	r.Route("/permissions", func(perms chi.Router) {
		perms.Group(func(admin chi.Router) {
			admin.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
			admin.Get("/", h.list)
			admin.Get("/{id}", h.getByID)
			admin.Get("/by-name/{name}", h.getByName)
			admin.Post("/", h.create)
			admin.Put("/{id}", h.update)
			admin.Delete("/{id}", h.softDelete)
			admin.Patch("/{id}/restore", h.restore)
		})
	})
}

func (h *PermissionsHandler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
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

func (h *PermissionsHandler) list(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	pagination := models.NewPagination(1, 50).WithMaxLimit(100)
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed >= 1 {
			pagination.Page = parsed
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed >= 1 {
			pagination.Limit = parsed
		}
	}

	total, err := h.service.Count(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	perms, err := h.service.List(r.Context(), pagination)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res := make([]*responses.PermissionResponse, len(perms))
	for i, p := range perms {
		res[i] = p.Response()
	}

	paginatedRes := responses.NewPaginatedResponse(res, pagination.GetPage(), pagination.GetLimit(), total)
	responses.WriteJSON(w, responses.RespondOk(paginatedRes, "Permissions retrieved"))
}

func (h *PermissionsHandler) getByID(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	perm, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			responses.WriteJSON(w, responses.RespondNotFound(nil, err.Error()))
			return
		}
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(), "Permission retrieved"))
}

func (h *PermissionsHandler) getByName(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "name is required"))
		return
	}

	perm, err := h.service.GetByName(r.Context(), name)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(), "Permission retrieved"))
}

func (h *PermissionsHandler) create(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	var dto dtos.CreatePermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Create(r.Context(), &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "Permission created"))
}

func (h *PermissionsHandler) update(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "id is required"))
		return
	}

	var dto dtos.UpdatePermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.ID = idParam

	if err := h.service.Update(r.Context(), &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Permission updated"))
}

func (h *PermissionsHandler) softDelete(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	if err := h.service.SoftDelete(r.Context(), id); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Permission deleted"))
}

func (h *PermissionsHandler) restore(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	if err := h.service.Restore(r.Context(), id); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Permission restored"))
}
