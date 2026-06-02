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
	service *service.PermissionService
}

func NewPermissionsHandler(svc *service.PermissionService) *PermissionsHandler {
	return &PermissionsHandler{
		service: svc,
	}
}

func (h *PermissionsHandler) Routes(r chi.Router) {
	r.Route("/permissions", func(p chi.Router) {
		p.Get("/", h.list)
		p.Route("/{id}", func(pid chi.Router) {
			p.Get("/", h.getByID)
			p.Put("/", h.update)
			p.Patch("/restore", h.restore)
			p.Delete("/", h.softDelete)
			p.Delete("/forever", h.delete)
		})
		p.Post("/", h.create)
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

	perms, err := middlewares.RequiresOwner(r, func(requester *models.User) bool {
		return requester.Role == models.ROLE_USER
	}, func(session *models.Session) (*[]*models.Permission, error) {
		dbPerms, err := h.service.List(r.Context(), pagination)
		if err != nil {
			return nil, err
		}

		return &dbPerms, nil
	})
	if err != nil {
		handleReqError(w, err)
		return
	}

	total, err := h.service.Count(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res := make([]*responses.PermissionResponse, len(*perms))
	for i, p := range *perms {
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

	perm, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.Permission, error) {
			perm, err := h.service.GetByID(r.Context(), pid)
			if err != nil {
				return nil, err
			}

			return perm, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(perm.Response(), "permission retrieved"))
}

func (h *PermissionsHandler) create(w http.ResponseWriter, r *http.Request) {
	_, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.Permission, error) {
			var dto dtos.CreatePermissionDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.service.Create(r.Context(), &dto); err != nil {
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

func (h *PermissionsHandler) update(w http.ResponseWriter, r *http.Request) {
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
		func(session *models.Session) (*models.Permission, error) {
			var dto dtos.UpdatePermissionDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.service.Update(r.Context(), pid, &dto); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.Permission, error) {
			if err := h.service.Restore(r.Context(), pid); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.Permission, error) {
			if err := h.service.SoftDelete(r.Context(), pid); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.Permission, error) {
			if err := h.service.Delete(r.Context(), pid); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission deleted"))
}
