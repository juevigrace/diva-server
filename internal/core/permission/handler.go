package permission

import (
	"net/http"
	"strconv"

	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type PermissionHandler struct {
	pRepo *PermissionRepo
}

func NewPermissionHandler(pRepo *PermissionRepo) *PermissionHandler {
	return &PermissionHandler{
		pRepo: pRepo,
	}
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

	perms, err := h.pRepo.List(r.Context(), pagination)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	total, err := h.pRepo.Count(r.Context())
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

	perm, err := h.pRepo.GetByID(r.Context(), pid)
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

	if err := h.pRepo.Create(r.Context(), &dto); err != nil {
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

	if err := h.pRepo.Update(r.Context(), pid, &dto); err != nil {
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

	if err := h.pRepo.Restore(r.Context(), pid); err != nil {
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

	if err := h.pRepo.SoftDelete(r.Context(), pid); err != nil {
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

	if err := h.pRepo.Delete(r.Context(), pid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "permission deleted"))
}
