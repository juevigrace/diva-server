package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserHandler struct {
	sService    *service.SessionService
	uService    *service.UserService
	upService   *service.UserPermissionService
	sHandler    *SessionHandler
	uaHandler   *UserActionsHandler
	upHandler   *UserPermissionHandler
	uprHandler  *UserPreferencesHandler
	uproHandler *UserProfileHandler
}

func NewUserHandler(
	sService *service.SessionService,
	uService *service.UserService,
	upService *service.UserPermissionService,
	sHandler *SessionHandler,
	uaHandler *UserActionsHandler,
	upHandler *UserPermissionHandler,
	uprHandler *UserPreferencesHandler,
	uproHandler *UserProfileHandler,
) *UserHandler {
	return &UserHandler{
		sService:    sService,
		uService:    uService,
		sHandler:    sHandler,
		uaHandler:   uaHandler,
		upHandler:   upHandler,
		uprHandler:  uprHandler,
		uproHandler: uproHandler,
		upService:   upService,
	}
}

func (h *UserHandler) Routes(r chi.Router) {
	r.Route("/user", func(u chi.Router) {
		u.Route("/check", func(check chi.Router) {
			check.Get("/username/{username}", h.checkUsername)
			check.Get("/email/{email}", h.checkEmail)
		})

		u.Group(func(auth chi.Router) {
			auth.Use(middlewares.SessionMiddleware(h.sService.GetByID))

			auth.Get("/", h.getAll)

			auth.Route("/{uid}", func(one chi.Router) {
				one.Get("/", h.getByID)

				// REQUIRES PERMISSION
				one.Patch("/email", h.updateEmail)
				one.Patch("/phone", h.updatePhone)
				one.Patch("/username", h.updateUsername)
				one.Patch("/password", h.updatePassword)
				one.Patch("/role", h.updateRole)
				one.Patch("/verified", h.updateVerified)

				one.Delete("/", h.softDelete)
				one.Delete("/forever", h.delete)

				h.sHandler.Routes(one)
				h.uprHandler.Routes(one)
				h.uproHandler.Routes(one)
			})

			auth.Post("/", h.create)

			h.uaHandler.Routes(auth)
			h.upHandler.Routes(auth)
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

	users, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*[]models.User, error) {
			dbUsers, err := h.uService.GetAll(r.Context(), pagination)
			if err != nil {
				return nil, err
			}

			return &dbUsers, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	total, err := h.uService.Count(r.Context())
	if err != nil {
		handleReqError(w, err)
		return
	}

	uRes := make([]*responses.UserResponse, len(*users))
	for i, u := range *users {
		uRes[i] = u.Response()
	}

	res := responses.NewPaginatedResponse(uRes, pagination.GetPage(), pagination.GetLimit(), total)

	responses.WriteJSON(w, responses.RespondOk(res, "users retrieved"))
}

func (h *UserHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != id
		},
		func(session *models.Session) (*models.User, error) {
			return h.uService.GetByID(r.Context(), id)
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(res.Response(), "user retrieved"))
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	_, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*int, error) {
			var dto dtos.CreateUserDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if _, err := h.uService.Create(r.Context(), &dto); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "user created"))
}

func (h *UserHandler) updateEmail(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != id
		},
		func(session *models.Session) (*any, error) {
			var dto dtos.UpdateEmailDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.uService.UpdateEmail(r.Context(), dto.Email, session.User.ID); err != nil {
				return nil, err
			}
			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "email updated"))
}

func (h *UserHandler) updatePassword(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (ret *any, err error) {
			deleteCall := func() {
				e := h.upService.Delete(r.Context(), session.User.ID, session.User.Permissions[models.PERMISSION_PASSWORD_UPDATE].Permission.ID)
				if e != nil && err == nil {
					err = e
				}
			}

			if err = middlewares.RequiresPermission(&session.User, models.PERMISSION_PASSWORD_UPDATE); err != nil {
				if errors.Is(err, models.ErrPermissionExpired) {
					deleteCall()
				}
				return
			}
			defer deleteCall()

			var dto dtos.UpdatePasswordDto
			if err = middlewares.ValidateBody(&dto, r); err != nil {
				return
			}

			if err = h.uService.UpdatePassword(r.Context(), uid, dto.NewPassword); err != nil {
				return
			}

			if err = h.sService.CloseAllByUser(r.Context(), uid); err != nil {
				return
			}

			return
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "password updated"))
}

func (h *UserHandler) updateUsername(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(r, func(requester *models.User) bool {
		return requester.Role == models.ROLE_USER && requester.ID != id
	}, func(session *models.Session) (*any, error) {
		var dto dtos.UpdateUsernameDto
		if err := middlewares.ValidateBody(&dto, r); err != nil {
			return nil, err
		}

		if err := h.uService.UpdateUsername(r.Context(), dto.Username, session.User.ID); err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "username updated"))
}

func (h *UserHandler) updatePhone(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(r, func(requester *models.User) bool {
		return requester.Role == models.ROLE_USER && requester.ID != id
	}, func(session *models.Session) (*any, error) {
		var dto dtos.UpdatePhoneNumberDto
		if err := middlewares.ValidateBody(&dto, r); err != nil {
			return nil, err
		}

		if err := h.uService.UpdatePhoneNumber(r.Context(), dto.PhoneNumber, id); err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "phone number updated"))
}

func (h *UserHandler) updateRole(w http.ResponseWriter, r *http.Request) {
	id, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*any, error) {
			var dto dtos.UpdateRole
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.uService.UpdateRole(r.Context(), models.RoleFromString(dto.Role), id); err != nil {
				return nil, err
			}

			return nil, err
		},
	)
	if err != nil {
		handleReqError(w, err)
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*any, error) {
			var dto dtos.UpdateVerified
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.uService.UpdateVerified(r.Context(), dto.Verified, id); err != nil {
				return nil, err
			}

			return nil, err
		},
	)
	if err != nil {
		handleReqError(w, err)
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != id
		},
		func(session *models.Session) (*any, error) {
			if err := h.uService.SoftDelete(r.Context(), id); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != id
		},
		func(session *models.Session) (*any, error) {
			if err := h.uService.Delete(r.Context(), id); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "user deleted forever"))
}
