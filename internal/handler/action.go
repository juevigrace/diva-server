package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type ActionHandler struct {
	service *service.ActionService
}

func NewActionHandler(svc *service.ActionService) *ActionHandler {
	return &ActionHandler{service: svc}
}

func (h *ActionHandler) Routes(r chi.Router) {
	r.Route("/actions", func(auth chi.Router) {
		auth.Get("/", h.getActions)
	})
}

func (h *ActionHandler) getActions(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	actions, err := h.service.GetAll(r.Context(), session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	result := make([]responses.ActionResponse, len(actions))
	for i, a := range actions {
		result[i] = responses.ActionResponse{ActionName: a.String()}
	}

	responses.WriteJSON(w, responses.RespondOk(result, "Success"))
}
