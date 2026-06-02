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

type UserPreferencesHandler struct {
	service *service.UserPreferencesService
}

func NewUserPreferencesHandler(svc *service.UserPreferencesService) *UserPreferencesHandler {
	return &UserPreferencesHandler{service: svc}
}

func (h *UserPreferencesHandler) Routes(r chi.Router) {
	r.Route("/preferences", func(pref chi.Router) {
		pref.Get("/", h.getByUser)
		pref.Route("/{pid}", func(uid chi.Router) {
			uid.Get("/", h.getByID)
			uid.Put("/", h.updatePreferences)
			uid.Delete("/", h.deletePreferences)
		})
		pref.Post("/", h.createPreferences)
	})
}

func (h *UserPreferencesHandler) getByUser(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	prefs, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*[]*models.UserPreferences, error) {
			prefs, err := h.service.GetByUser(r.Context(), uid)
			if err != nil {
				return nil, err
			}

			return &prefs, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	res := make([]*responses.UserPreferencesResponse, len(*prefs))
	for i, p := range *prefs {
		res[i] = p.Response(&uid)
	}

	responses.WriteJSON(w, responses.RespondOk(res, "preferences retrieved"))
}

func (h *UserPreferencesHandler) getByID(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	pref, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*models.UserPreferences, error) {
			prefs, err := h.service.GetByID(r.Context(), pid)
			if err != nil {
				return nil, err
			}

			return prefs, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(pref.Response(&uid), "preference retrieved"))
}

func (h *UserPreferencesHandler) createPreferences(w http.ResponseWriter, r *http.Request) {
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
		func(session *models.Session) (*any, error) {
			var dto dtos.CreateUserPreferencesDto
			if err := middlewares.ValidateBody(&dto, r); err != nil {
				return nil, err
			}

			dto.Device = session.Device

			if err := h.service.Create(r.Context(), session.User.ID, &dto); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(nil, "preferences created"))
}

func (h *UserPreferencesHandler) updatePreferences(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*any, error) {
			var dto dtos.UpdateUserPreferencesDto
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

	responses.WriteJSON(w, responses.RespondAccepted(nil, "preferences updated"))
}

func (h *UserPreferencesHandler) deletePreferences(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	pid, err := middlewares.GetUUIDFromURL(r, "pid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*any, error) {
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

	responses.WriteJSON(w, responses.RespondNoContent(nil, "preference deleted"))
}
