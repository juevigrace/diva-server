package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	r.Route("/actions", func(act chi.Router) {
		act.Group(func(rg chi.Router) {
			rg.Use(middlewares.RequireResourceOwner(
				"uid",
				func(_ context.Context, reqid, resid uuid.UUID) (any, bool) {
					return nil, reqid == resid
				},
				models.PERMISSION_ACTIONS_READ,
			))
			rg.Get("/", h.getAll)
		})
	})
}

func (h *UserActionsHandler) Routes(r chi.Router) {
	r.Route("/actions", func(ac chi.Router) {
		ac.Route("/{aid}", func(aid chi.Router) {
			aid.With(middlewares.RequireResourceOwner(
				"aid",
				func(ctx context.Context, reqid, resid uuid.UUID) (any, bool) {
					action, err := h.service.GetOneByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if action.UserID != reqid {
						return nil, false
					}
					return action, true
				},
				models.PERMISSION_ACTIONS_READ,
			)).Get("/", h.getByID)
			aid.With(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
				middlewares.RequirePermission(models.PERMISSION_ACTIONS_WRITE),
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

	res := make([]*responses.UserActionResponse, len(actions))
	for i, a := range actions {
		res[i] = a.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(res, "actions retrieved"))
}

func (h *UserActionsHandler) getByID(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	action, ok := rc.Cache["aid"].(*models.UserAction)
	if !ok {
		actionID, err := middlewares.GetUUIDFromURL(r, "aid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
		action, err = h.service.GetOneByID(r.Context(), actionID)
		if err != nil {
			responses.HandleReqError(w, err)
			return
		}
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
