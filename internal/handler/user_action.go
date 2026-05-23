package handler

import (
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

func (h *UserActionsHandler) Routes(r chi.Router) {
	r.Route("/actions", func(auth chi.Router) {
		auth.Get("/", h.getActions)

		auth.Group(func(admin chi.Router) {
			admin.Use(middlewares.SessionMiddleware(h.sessionService.GetByID))
			admin.Get("/{id}", h.getActionByID)
			admin.Delete("/{id}", h.deleteAction)
		})
	})
}

func (h *UserActionsHandler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return false
	}
	if session.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "admin access required"))
		return false
	}
	return true
}

func (h *UserActionsHandler) getActions(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	actions, err := h.service.GetAllByUser(r.Context(), session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	result := make([]*responses.UserActionResponse, len(actions))
	for i, a := range actions {
		result[i] = a.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(result, "Success"))
}

func (h *UserActionsHandler) getActionByID(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	action, err := h.service.GetOneByID(r.Context(), id)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(action.Response(), "Action retrieved"))
}

func (h *UserActionsHandler) deleteAction(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Action deleted"))
}
