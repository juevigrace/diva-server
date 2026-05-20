package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserActionsHandler struct {
	service *service.UserActionsService
}

func NewUserActionsHandler(svc *service.UserActionsService) *UserActionsHandler {
	return &UserActionsHandler{service: svc}
}

func (h *UserActionsHandler) Routes(r chi.Router) {
	r.Route("/actions", func(auth chi.Router) {
		auth.Get("/", h.getActions)
	})
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
