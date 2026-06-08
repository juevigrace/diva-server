package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type PermissionsHandler struct {
	service  *service.PermissionService
	sService *service.SessionService
}

func NewPermissionsHandler(
	svc *service.PermissionService,
	sService *service.SessionService,
) *PermissionsHandler {
	return &PermissionsHandler{
		service:  svc,
		sService: sService,
	}
}

func (h *PermissionsHandler) Routes(r chi.Router) {
	r.Route("/permissions", func(p chi.Router) {
		p.Use(middlewares.RequiresSession(h.sService.GetByID))

		p.With(
			middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_READ),
		).Get("/", h.list)

		p.Route("/{pid}", func(pid chi.Router) {
			pid.With(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
				middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_READ),
			).Get("/", h.getByID)

			pid.Group(func(wg chi.Router) {
				wg.Use(
					middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
					middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_WRITE),
				)
				wg.Put("/", h.update)
				wg.Patch("/restore", h.restore)
				wg.Delete("/", h.softDelete)
				wg.Delete("/forever", h.delete)
			})
		})

		p.With(
			middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_WRITE),
		).Post("/", h.create)
	})
}

func (h *PermissionsHandler) list(w http.ResponseWriter, r *http.Request) {
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

	perms, err := h.service.List(r.Context(), pagination)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	total, err := h.service.Count(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res := make([]*responses.PermissionResponse, len(perms))
	for i, p := range perms {
		res[i] = p.Response()
	}

	paginatedRes := responses.NewPaginatedResponse(res, pagination.GetPage(), pagination.GetLimit(), total)
	responses.WriteJSON(w, responses.RespondOk(paginatedRes, "permissions retrieved"))
}

func (h *PermissionsHandler) getByID(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	perm, err := h.service.GetByID(r.Context(), pid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(), "permission retrieved"))
}

func (h *PermissionsHandler) create(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreatePermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Create(r.Context(), &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "permission created"))
}

func (h *PermissionsHandler) update(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	var dto dtos.UpdatePermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Update(r.Context(), pid, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission updated"))
}

func (h *PermissionsHandler) restore(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Restore(r.Context(), pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission restored"))
}

func (h *PermissionsHandler) softDelete(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.SoftDelete(r.Context(), pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission deleted"))
}

func (h *PermissionsHandler) delete(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.service.Delete(r.Context(), pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission deleted"))
}
