package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserProfileHandler struct {
	upService   *service.UserPermissionService
	uproService *service.UserProfileService
}

func NewUserProfileHandler(
	upService *service.UserPermissionService,
	uproService *service.UserProfileService,
) *UserProfileHandler {
	return &UserProfileHandler{
		upService:   upService,
		uproService: uproService,
	}
}

func (h *UserProfileHandler) UserRoutes(r chi.Router) {
	r.Route("/profile", func(pr chi.Router) {
		pr.Get("/", h.getOne)
		pr.With(
			middlewares.RequirePermission(models.PERMISSION_USERS_PROFILE_WRITE),
			middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"uid"},
					Perms:     []models.PermissionAction{models.PERMISSION_USERS_PROFILE_WRITE},
				},
				func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					return nil, reqid == resid
				},
			),
		).Post("/", h.create)
		pr.Group(func(rg chi.Router) {
			rg.Use(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"uid"},
					Perms:     []models.PermissionAction{models.PERMISSION_USERS_PROFILE_WRITE},
				},
				func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					return nil, reqid == resid
				},
			))
			rg.Put("/", h.update)
			rg.Patch("/avatar", h.updateAvatar)
		})
	})
}

func (h *UserProfileHandler) getOne(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	dbProfile, err := h.uproService.GetByUserID(r.Context(), uid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(dbProfile.Response(), "profile retrieved"))
}

func (h *UserProfileHandler) create(w http.ResponseWriter, r *http.Request) {
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

	var dto dtos.CreateProfileDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		return
	}

	if err = h.uproService.Create(r.Context(), uid, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.ID == uid {
		if perm, ok := session.User.Permissions[models.PERMISSION_USERS_PROFILE_WRITE]; ok {
			if err := h.upService.Delete(r.Context(), session.User.ID, perm.Permission.ID); err != nil {
				responses.HandleReqError(w, err)
				return
			}
		}
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "profile created"))
}

func (h *UserProfileHandler) update(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	var dto dtos.UpdateProfileDto
	if err = middlewares.ValidateBody(&dto, r); err != nil {
		return
	}

	if err = h.uproService.Update(r.Context(), uid, &dto); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "profile updated"))
}

func (h *UserProfileHandler) updateAvatar(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err = h.uproService.UpdateAvatar(r.Context(), uid, ""); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "avatar updated"))
}
