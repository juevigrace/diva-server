package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type SessionHandler struct {
	sService *service.SessionService
}

func NewSessionHandler(sService *service.SessionService) *SessionHandler {
	return &SessionHandler{sService: sService}
}

func (h *SessionHandler) UserRoutes(r chi.Router) {
	r.Route("/sessions", func(s chi.Router) {
		s.Get("/", h.listByUser)
		s.Route("/{sid}", func(sid chi.Router) {
			sid.Get("/", h.getByID)
			sid.Delete("/", h.close)
		})
	})
}

func (h *SessionHandler) Routes(r chi.Router) {
	r.Route("/sessions", func(s chi.Router) {
		s.Use(middlewares.RequiresSession(h.sService.GetByID))
		s.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
		s.Delete("/expired", h.deleteExpired)
	})
}

func (h *SessionHandler) listByUser(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	sessions, err := h.sService.GetByUser(r.Context(), uid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	res := make([]*responses.SessionResponse, len(sessions))
	for i, s := range sessions {
		res[i] = s.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(res, "sessions retrieved"))
}

func (h *SessionHandler) getByID(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	sid, err := middlewares.GetUUIDFromURL(r, "sid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	dbSession, err := h.sService.GetByID(r.Context(), sid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if dbSession.User.ID != uid {
		responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrForbidden.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(dbSession.Response(), "session retrieved"))
}

func (h *SessionHandler) close(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	sid, err := middlewares.GetUUIDFromURL(r, "sid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, errs.ErrSessionNotFound.Error()))
		return
	}

	target, err := h.sService.GetByID(r.Context(), sid)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	if session.User.Role == models.ROLE_USER && target.User.ID != uid {
		responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrForbidden.Error()))
		return
	}

	if err := h.sService.Close(r.Context(), sid); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "session closed"))
}

func (h *SessionHandler) deleteExpired(w http.ResponseWriter, r *http.Request) {
	if err := h.sService.DeleteExpired(r.Context()); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "expired sessions deleted"))
}
