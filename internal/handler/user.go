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

type UserHandler struct {
	sessionService  *service.SessionService
	userService     *service.UserService
	userMeHandler   *UserMeHandler
	userPermHandler *UserPermissionHandler
}

func NewUserHandler(
	sessionService *service.SessionService,
	userService *service.UserService,
	userMeHandler *UserMeHandler,
	userPermHandler *UserPermissionHandler,
) *UserHandler {
	return &UserHandler{
		sessionService:  sessionService,
		userService:     userService,
		userMeHandler:   userMeHandler,
		userPermHandler: userPermHandler,
	}
}

func (h *UserHandler) Routes(r chi.Router) {
	r.Route("/user", func(u chi.Router) {
		u.Route("/check", func(check chi.Router) {
			check.Get("/username/{username}", h.checkUsername)
			check.Get("/email/{email}", h.checkEmail)
		})

		u.Route("/{id}", func(uid chi.Router) {
			uid.Get("/", h.getUserByID)
			uid.Group(func(admin chi.Router) {
				admin.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
				admin.Put("/", h.updateUser)
				admin.Delete("/", h.deleteUser)
			})
		})

		u.Group(func(admin chi.Router) {
			admin.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
			admin.Get("/", h.getUsers)
			admin.Post("/", h.createUser)
		})

		u.Group(func(auth chi.Router) {
			auth.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
			h.userMeHandler.Routes(auth)
			h.userPermHandler.Routes(auth)
		})
	})
}

func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	if session.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "admin access required"))
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

	total, err := h.userService.Count(r.Context())
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	users, err := h.userService.GetAll(r.Context(), pagination)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res := make([]*responses.UserResponse, len(users))
	for i, u := range users {
		res[i] = &responses.UserResponse{
			ID:           u.ID.String(),
			Email:        u.Email,
			Username:     u.Username,
			BirthDate:    u.BirthDate,
			PhoneNumber:  u.PhoneNumber,
			Alias:        u.Alias,
			Avatar:       u.Avatar,
			Bio:          u.Bio,
			UserVerified: u.UserVerified,
			Role:         u.Role.String(),
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
			DeletedAt:    u.DeletedAt,
		}
	}

	paginatedRes := responses.NewPaginatedResponse(res, pagination.GetPage(), pagination.GetLimit(), total)
	responses.WriteJSON(w, responses.RespondOk(paginatedRes, "Users retrieved"))
}

func (h *UserHandler) checkUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "username is required"))
		return
	}

	available, err := h.userService.CheckUsernameAvailable(r.Context(), username)
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

	available, err := h.userService.CheckEmailAvailable(r.Context(), email)
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

func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "id is required"))
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			responses.WriteJSON(w, responses.RespondNotFound(nil, err.Error()))
			return
		}
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(&responses.UserResponse{
		ID:           user.ID.String(),
		Email:        "",
		Username:     user.Username,
		BirthDate:    user.BirthDate,
		PhoneNumber:  "",
		Alias:        user.Alias,
		Avatar:       user.Avatar,
		Bio:          user.Bio,
		UserVerified: user.UserVerified,
		Role:         "",
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		DeletedAt:    user.DeletedAt,
	}, "User retrieved"))
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	if session.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "admin access required"))
		return
	}

	var dto dtos.CreateUserDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	userID, err := h.userService.Create(r.Context(), &dto)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(map[string]string{"id": userID.String()}, "User created"))
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	if session.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "admin access required"))
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	var dto dtos.UpdateProfileDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.userService.UpdateProfile(r.Context(), id, &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "User updated"))
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	if session.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "admin access required"))
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "User deleted"))
}
