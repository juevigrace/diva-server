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

type SessionHandler struct {
	Service *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{Service: svc}
}

func (h *SessionHandler) Routes(r chi.Router) {
	r.Route("/sessions", func(sessions chi.Router) {
		sessions.Group(func(protected chi.Router) {
			protected.Use(middlewares.SessionMiddleware(h.Service.GetByID))
			protected.Get("/", h.list)
			protected.Get("/{id}", h.getByID)
			protected.Delete("/{id}", h.close)
		})

		sessions.Group(func(admin chi.Router) {
			admin.Use(middlewares.SessionMiddleware(h.Service.GetByID))
			admin.Delete("/expired", h.deleteExpired)
		})
	})
}

func (h *SessionHandler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
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

func (h *SessionHandler) list(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	sessions, err := h.Service.GetByUser(r.Context(), session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res := make([]*responses.SessionResponse, len(sessions))
	for i, s := range sessions {
		res[i] = models.ToSessionResponse(s)
	}

	responses.WriteJSON(w, responses.RespondOk(res, "Sessions retrieved"))
}

func (h *SessionHandler) getByID(w http.ResponseWriter, r *http.Request) {
	currentSession, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	s, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if currentSession.User.ID != s.User.ID && currentSession.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "access denied"))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(models.ToSessionResponse(s), "Session retrieved"))
}

func (h *SessionHandler) close(w http.ResponseWriter, r *http.Request) {
	currentSession, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid id format"))
		return
	}

	target, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if currentSession.User.ID != target.User.ID && currentSession.User.Role != models.ROLE_ADMIN {
		responses.WriteJSON(w, responses.RespondForbidden(nil, "access denied"))
		return
	}

	if err := h.Service.Close(r.Context(), id); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Session closed"))
}

func (h *SessionHandler) deleteExpired(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	if err := h.Service.DeleteExpired(r.Context()); err != nil {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Expired sessions deleted"))
}
