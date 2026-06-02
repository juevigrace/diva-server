package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserActionsHandler struct {
	service        *service.UserActionsService
	sessionService *service.SessionService
}

func NewUserActionsHandler(svc *service.UserActionsService, sessionService *service.SessionService) *UserActionsHandler {
	return &UserActionsHandler{service: svc, sessionService: sessionService}
}

func (h *UserActionsHandler) Routes(r chi.Router) {
	r.Route("/{uid}/actions", func(uid chi.Router) {
		uid.Get("/", h.getAll)
		uid.Get("/{aid}", h.getByID)
	})
	r.Route("/actions", func(ac chi.Router) {
		ac.Route("/{aid}", func(aid chi.Router) {
			aid.Delete("/", h.deleteAction)
		})
	})
}

func (h *UserActionsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	actions, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*[]models.UserAction, error) {
			dbActions, err := h.service.GetAllByUser(r.Context(), uid)
			if err != nil {
				return nil, err
			}

			return &dbActions, err
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	res := make([]*responses.UserActionResponse, len(*actions))
	for i, a := range *actions {
		res[i] = a.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(res, "actions retrieved"))
}

func (h *UserActionsHandler) getByID(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	actionID, err := middlewares.GetUUIDFromURL(r, "aid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	action, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*models.UserAction, error) {
			dbAction, err := h.service.GetOneByID(r.Context(), actionID)
			if err != nil {
				return nil, err
			}

			if dbAction.UserID != uid {
				return nil, models.ErrForbidden
			}

			return dbAction, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(action.Response(), "action retrieved"))
}

func (h *UserActionsHandler) deleteAction(w http.ResponseWriter, r *http.Request) {
	actionID, err := middlewares.GetUUIDFromURL(r, "aid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER
		},
		func(session *models.Session) (*models.UserAction, error) {
			if err := h.service.Delete(r.Context(), actionID); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "action deleted"))
}
