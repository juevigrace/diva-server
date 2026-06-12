package permission

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/core"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type PermissionHandler struct {
	pService *PermissionService
	provider core.Provider[*models.Session]
}

func NewPermissionHandler(
	pService *PermissionService,
	provider core.Provider[*models.Session],
) *PermissionHandler {
	return &PermissionHandler{
		pService: pService,
		provider: provider,
	}
}

func (h *PermissionHandler) Routes(r chi.Router) {
	r.Route("/permissions", func(p chi.Router) {
		p.Use(middlewares.RequiresSession(h.provider.GetByID))

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
			})
			pid.Group(func(wg chi.Router) {
				wg.Use(middlewares.RequireRole(models.ROLE_ADMIN))
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

func (h *PermissionHandler) list(w http.ResponseWriter, r *http.Request) {
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

	perms, err := h.pService.List(r.Context(), pagination)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	total, err := h.pService.Count(r.Context())
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

func (h *PermissionHandler) getByID(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	perm, err := h.pService.GetByID(r.Context(), pid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(), "permission retrieved"))
}

func (h *PermissionHandler) create(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreatePermissionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
		return
	}

	requestedLevel := models.RoleFromString(dto.RoleLevel)
	if session.User.Role < requestedLevel {
		responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrForbidden.Error()))
		return
	}

	if err := h.pService.Create(r.Context(), &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "permission created"))
}

func (h *PermissionHandler) update(w http.ResponseWriter, r *http.Request) {
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

	if err := h.pService.Update(r.Context(), pid, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission updated"))
}

func (h *PermissionHandler) restore(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.pService.Restore(r.Context(), pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission restored"))
}

func (h *PermissionHandler) softDelete(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.pService.SoftDelete(r.Context(), pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission deleted"))
}

func (h *PermissionHandler) delete(w http.ResponseWriter, r *http.Request) {
	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.pService.Delete(r.Context(), pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission deleted"))
}
