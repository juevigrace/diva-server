package actions

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/core/session"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

type UserActionsHandler struct {
	uaService *UserActionsService
	sService  *session.SessionService
}

func NewUserActionsHandler(
	uaService *UserActionsService,
	sService *session.SessionService,
) *UserActionsHandler {
	return &UserActionsHandler{
		uaService: uaService,
		sService:  sService,
	}
}

func (h *UserActionsHandler) UserRoutes(r chi.Router) {
	r.Route("/actions", func(act chi.Router) {
		act.Group(func(rg chi.Router) {
			rg.Use(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"uid"},
					Perms:     []models.PermissionAction{models.PERMISSION_ACTIONS_READ},
				},
				func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					if reqid != resid {
						return nil, false
					}
					return map[string]any{"uid": resid}, true
				},
			))
			rg.Get("/", h.getAll)
		})
	})
}

func (h *UserActionsHandler) Routes(r chi.Router) {
	r.Route("/actions", func(ac chi.Router) {
		ac.Route("/{aid}", func(aid chi.Router) {
			aid.With(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"aid"},
					Perms:     []models.PermissionAction{models.PERMISSION_ACTIONS_READ},
				},
				func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					action, err := h.uaService.GetOneByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if action.UserID != reqid {
						return nil, false
					}
					return map[string]any{"aid": action}, true
				},
			)).Get("/", h.getByID)
			aid.With(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
				middlewares.RequirePermission(models.PERMISSION_ACTIONS_WRITE),
			).Delete("/", h.deleteAction)
		})
	})
}

func (h *UserActionsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
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

	actions, err := h.uaService.GetAllByUser(r.Context(), uid)
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
		action, err = h.uaService.GetOneByID(r.Context(), actionID)
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

	if err := h.uaService.Delete(r.Context(), actionID); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondNoContent(nil, "action deleted"))
}
