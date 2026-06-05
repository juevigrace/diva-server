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

func (h *UserActionsHandler) UserRoutes(r chi.Router) {
	r.Route("/actions", func(a chi.Router) {
		a.Get("/", h.getAll)
		a.Get("/{aid}", h.getByID)
	})
}

func (h *UserActionsHandler) Routes(r chi.Router) {
	r.Route("/actions", func(ac chi.Router) {
		ac.Route("/{aid}", func(aid chi.Router) {
			aid.With(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			).Delete("/", h.deleteAction)
		})
	})
}

func (h *UserActionsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	actions, err := h.service.GetAllByUser(r.Context(), uid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	res := make([]*responses.UserActionResponse, len(actions))
	for i, a := range actions {
		res[i] = a.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(res, "actions retrieved"))
}

func (h *UserActionsHandler) getByID(w http.ResponseWriter, r *http.Request) {
	actionID, err := middlewares.GetUUIDFromURL(r, "aid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	action, err := h.service.GetOneByID(r.Context(), actionID)
	if err != nil {
		responses.HandleReqError(w, err)
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

	if err := h.service.Delete(r.Context(), actionID); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "action deleted"))
}
