package user

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type UserHandler struct {
	uRepo  *UserRepo
	usRepo *UserStateRepo
}

func NewUserHandler(
	uRepo *UserRepo,
	usRepo *UserStateRepo,
) *UserHandler {
	return &UserHandler{
		uRepo:  uRepo,
		usRepo: usRepo,
	}
}

func (h *UserHandler) checkUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "username is required"))
		return
	}

	available, err := h.uRepo.CheckUsernameAvailable(r.Context(), username)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if !available {
		responses.WriteJSON(w, responses.RespondConflict(nil, "username already taken"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "username is available"))
}

func (h *UserHandler) checkEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "email is required"))
		return
	}

	available, err := h.uRepo.CheckEmailAvailable(r.Context(), email)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if !available {
		responses.WriteJSON(w, responses.RespondConflict(nil, "email already taken"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "email is available"))
}

func (h *UserHandler) getAll(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.uRepo.GetAll(r.Context(), pagination)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	total, err := h.uRepo.Count(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	uRes := make([]*responses.UserResponse, len(users))
	for i, u := range users {
		uRes[i] = u.Response()
	}

	res := responses.NewPaginatedResponse(uRes, pagination.GetPage(), pagination.GetLimit(), total)

	responses.WriteJSON(w, responses.RespondOk(res, "users retrieved"))
}

func (h *UserHandler) getByID(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
		return
	}

	dbUser, err := h.uRepo.GetByID(r.Context(), uid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.Role == models.ROLE_USER {
		dbUser.Email = ""
		dbUser.PhoneNumber = ""
	}

	responses.WriteJSON(w, responses.RespondOk(dbUser.Response(), "user retrieved"))
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreateUserDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if _, err := h.uRepo.Create(r.Context(), &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "user created"))
}

func (h *UserHandler) updateEmail(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
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

	var dto dtos.UpdateEmailDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uRepo.UpdateEmail(r.Context(), rc.Session, dto.Email, uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "email updated"))
}

func (h *UserHandler) updatePassword(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
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

	var dto dtos.UpdatePasswordDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uRepo.UpdatePassword(r.Context(), rc.Session, uid, dto.NewPassword); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "password updated"))
}

func (h *UserHandler) updateUsername(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
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

	var dto dtos.UpdateUsernameDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uRepo.UpdateUsername(r.Context(), rc.Session, dto.Username, uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "username updated"))
}

func (h *UserHandler) updatePhone(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
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

	var dto dtos.UpdatePhoneNumberDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uRepo.UpdatePhoneNumber(r.Context(), rc.Session, dto.PhoneNumber, uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "phone number updated"))
}

func (h *UserHandler) updateRole(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	var dto dtos.UpdateRole
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	targetRole := models.RoleFromString(dto.Role)
	if rc.Session.User.Role < targetRole {
		responses.WriteJSON(w, responses.RespondForbidden(nil, errs.ErrPermissionDenied.Error()))
		return
	}

	if err = h.uRepo.UpdateRole(r.Context(), targetRole, id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "role updated"))
}

func (h *UserHandler) updateVerified(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	var dto dtos.UpdateVerified
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.usRepo.UpdateVerified(r.Context(), dto.Verified, id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "verified updated"))
}

func (h *UserHandler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	var dto dtos.UpdateUserStatus
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.usRepo.UpdateStatus(r.Context(), models.UserStatusFromString(dto.Status), id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "status updated"))
}

func (h *UserHandler) pingStatus(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
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

	if err := h.usRepo.UpdateLastActiveAt(r.Context(), uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "status pinged"))
}

func (h *UserHandler) softDelete(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
		return
	}

	id, ok := rc.Cache["uid"].(uuid.UUID)
	if !ok {
		id, err = middlewares.GetUUIDFromURL(r, "uid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	}

	if err := h.uRepo.SoftDelete(r.Context(), id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "user deleted"))
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
		return
	}

	id, ok := rc.Cache["uid"].(uuid.UUID)
	if !ok {
		id, err = middlewares.GetUUIDFromURL(r, "uid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	}

	if err := h.uRepo.Delete(r.Context(), id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "user deleted forever"))
}

func (h *UserHandler) restore(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.uRepo.Restore(r.Context(), id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "user restored"))
}
