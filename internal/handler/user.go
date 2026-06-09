package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserHandler struct {
	uService    *service.UserService
	sHandler    *SessionHandler
	uaHandler   *UserActionsHandler
	upHandler   *UserPermissionHandler
	uprHandler  *UserPreferencesHandler
	uproHandler *UserProfileHandler
}

func NewUserHandler(
	uService *service.UserService,
	sHandler *SessionHandler,
	uaHandler *UserActionsHandler,
	upHandler *UserPermissionHandler,
	uprHandler *UserPreferencesHandler,
	uproHandler *UserProfileHandler,
) *UserHandler {
	return &UserHandler{
		uService:    uService,
		uaHandler:   uaHandler,
		upHandler:   upHandler,
		uprHandler:  uprHandler,
		uproHandler: uproHandler,
		sHandler:    sHandler,
	}
}

func (h *UserHandler) Routes(r chi.Router) {
	r.Route("/user", func(u chi.Router) {
		u.Route("/check", func(check chi.Router) {
			check.Get("/username/{username}", h.checkUsername)
			check.Get("/email/{email}", h.checkEmail)
		})

		u.Group(func(auth chi.Router) {
			auth.Use(middlewares.RequiresSession(h.sHandler.sService.GetByID))

			auth.Group(func(admin chi.Router) {
				admin.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
				admin.Get("/", h.getAll)
				admin.Post("/", h.create)
			})

			auth.Route("/{uid}", func(uid chi.Router) {
				uid.Get("/", h.getByID)

			uid.With(
				middlewares.RequirePermission(models.PERMISSION_USERS_EMAIL_WRITE),
				middlewares.RequireResourceOwner(
					&middlewares.RequireOwnerParams{
						UrlParams: []string{"uid"},
						Perms:     []models.PermissionAction{models.PERMISSION_USERS_EMAIL_WRITE},
					},
					func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
						resid, err := uuid.Parse(resParams[0])
						if err != nil {
							return nil, false
						}
						return nil, reqid == resid
					},
				),
			).Patch("/email", h.updateEmail)
			uid.With(
				middlewares.RequirePermission(models.PERMISSION_USERS_PHONE_WRITE),
				middlewares.RequireResourceOwner(
					&middlewares.RequireOwnerParams{
						UrlParams: []string{"uid"},
						Perms:     []models.PermissionAction{models.PERMISSION_USERS_PHONE_WRITE},
					},
					func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
						resid, err := uuid.Parse(resParams[0])
						if err != nil {
							return nil, false
						}
						return nil, reqid == resid
					},
				),
			).Patch("/phone", h.updatePhone)
			uid.With(
				middlewares.RequirePermission(models.PERMISSION_USERS_USERNAME_WRITE),
				middlewares.RequireResourceOwner(
					&middlewares.RequireOwnerParams{
						UrlParams: []string{"uid"},
						Perms:     []models.PermissionAction{models.PERMISSION_USERS_USERNAME_WRITE},
					},
					func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
						resid, err := uuid.Parse(resParams[0])
						if err != nil {
							return nil, false
						}
						return nil, reqid == resid
					},
				),
			).Patch("/username", h.updateUsername)
			uid.With(
				middlewares.RequirePermission(models.PERMISSION_USERS_PASSWORD_WRITE),
				middlewares.RequireResourceOwner(
					&middlewares.RequireOwnerParams{
						UrlParams: []string{"uid"},
						Perms:     []models.PermissionAction{models.PERMISSION_USERS_PASSWORD_WRITE},
					},
					func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
						resid, err := uuid.Parse(resParams[0])
						if err != nil {
							return nil, false
						}
						return nil, reqid == resid
					},
				),
			).Patch("/password", h.updatePassword)

				uid.Group(func(admin chi.Router) {
					admin.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
					admin.With(middlewares.RequirePermission(models.PERMISSION_USERS_ROLE_WRITE)).Patch("/role", h.updateRole)
					admin.With(middlewares.RequirePermission(models.PERMISSION_USERS_VERIFIED_WRITE)).Patch("/verified", h.updateVerified)
				})

			uid.Group(func(wg chi.Router) {
				wg.Use(middlewares.RequireResourceOwner(
					&middlewares.RequireOwnerParams{
						UrlParams: []string{"uid"},
						Perms:     []models.PermissionAction{models.PERMISSION_USERS_WRITE},
					},
					func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
						resid, err := uuid.Parse(resParams[0])
						if err != nil {
							return nil, false
						}
						return nil, reqid == resid
					},
				))
				wg.Delete("/", h.softDelete)
				wg.Delete("/forever", h.delete)
			})

				h.uaHandler.UserRoutes(uid)
				h.upHandler.UserRoutes(uid)
				h.uprHandler.UserRoutes(uid)
				h.uproHandler.UserRoutes(uid)
				h.sHandler.UserRoutes(uid)
			})

			h.uaHandler.Routes(auth)
			h.uprHandler.Routes(auth)
		})
	})
}

func (h *UserHandler) checkUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "username is required"))
		return
	}

	available, err := h.uService.CheckUsernameAvailable(r.Context(), username)
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

	available, err := h.uService.CheckEmailAvailable(r.Context(), email)
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

	users, err := h.uService.GetAll(r.Context(), pagination)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	total, err := h.uService.Count(r.Context())
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

	dbUser, err := h.uService.GetByID(r.Context(), uid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	// TODO: create a permission for this?
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

	if _, err := h.uService.Create(r.Context(), &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "user created"))
}

func (h *UserHandler) updateEmail(w http.ResponseWriter, r *http.Request) {
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

	var dto dtos.UpdateEmailDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uService.UpdateEmail(r.Context(), dto.Email, uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.ID == uid {
		if err := h.upHandler.service.Delete(r.Context(), session.User.ID, session.User.Permissions[models.PERMISSION_USERS_EMAIL_WRITE].Permission.ID); err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "email updated"))
}

func (h *UserHandler) updatePassword(w http.ResponseWriter, r *http.Request) {
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

	var dto dtos.UpdatePasswordDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uService.UpdatePassword(r.Context(), uid, dto.NewPassword); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if err = h.sHandler.sService.CloseAllByUser(r.Context(), uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.ID == uid {
		if err := h.upHandler.service.Delete(r.Context(), session.User.ID, session.User.Permissions[models.PERMISSION_USERS_PASSWORD_WRITE].Permission.ID); err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "password updated"))
}

func (h *UserHandler) updateUsername(w http.ResponseWriter, r *http.Request) {
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

	var dto dtos.UpdateUsernameDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uService.UpdateUsername(r.Context(), dto.Username, uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.ID == uid {
		if err := h.upHandler.service.Delete(r.Context(), session.User.ID, session.User.Permissions[models.PERMISSION_USERS_USERNAME_WRITE].Permission.ID); err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "username updated"))
}

func (h *UserHandler) updatePhone(w http.ResponseWriter, r *http.Request) {
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

	var dto dtos.UpdatePhoneNumberDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uService.UpdatePhoneNumber(r.Context(), dto.PhoneNumber, uid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.ID == uid {
		if err := h.upHandler.service.Delete(r.Context(), session.User.ID, session.User.Permissions[models.PERMISSION_USERS_PHONE_WRITE].Permission.ID); err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "phone number updated"))
}

func (h *UserHandler) updateRole(w http.ResponseWriter, r *http.Request) {
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

	if err = h.uService.UpdateRole(r.Context(), models.RoleFromString(dto.Role), id); err != nil {
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

	if err = h.uService.UpdateVerified(r.Context(), dto.Verified, id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "verified updated"))
}

func (h *UserHandler) softDelete(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.uService.SoftDelete(r.Context(), id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "user deleted"))
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.uService.Delete(r.Context(), id); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "user deleted forever"))
}
