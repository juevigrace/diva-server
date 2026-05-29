package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserProfileHandler struct {
	upService *service.UserProfileService
}

func NewUserProfileHandler(
	upService *service.UserProfileService,
) *UserProfileHandler {
	return &UserProfileHandler{
		upService: upService,
	}
}

func (h *UserProfileHandler) Routes(r chi.Router) {
	r.Route("/profile", func(pr chi.Router) {
		pr.Get("/", h.getOne)
		pr.Post("/", h.create)
		pr.Put("/", h.update)
		pr.Patch("/avatar", h.updateAvatar)
	})
}

func (h *UserProfileHandler) getOne(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	profile, err := middlewares.RequiresOwnerOrPerms(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*models.UserProfile, error) {
			dbProfile, err := h.upService.GetByUserID(r.Context(), uid)
			if err != nil {
				return nil, err
			}

			return dbProfile, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(profile.Response(&uid), "profile retrieved"))
}

func (h *UserProfileHandler) create(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwnerOrPerms(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*any, error) {
			var dto dtos.CreateProfileDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.upService.Create(r.Context(), uid, &dto); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "profile created"))
}

func (h *UserProfileHandler) update(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwnerOrPerms(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*any, error) {
			var dto dtos.UpdateProfileDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			if err := h.upService.Update(r.Context(), uid, &dto); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
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

	_, err = middlewares.RequiresOwnerOrPerms(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*any, error) {
			if err := h.upService.UpdateAvatar(r.Context(), uid, ""); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondAccepted(nil, "avatar updated"))
}
