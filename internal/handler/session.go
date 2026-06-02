package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	r.Route("/sessions", func(s chi.Router) {
		s.Get("/", h.list)
		s.Get("/{sid}", h.getByID)
		s.Delete("/{sid}", h.close)
		s.Delete("/expired", h.deleteExpired)
	})
}

func (h *SessionHandler) list(w http.ResponseWriter, r *http.Request) {
	uid, err := middlewares.GetUUIDFromURL(r, "uid")
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	sessions, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*[]*models.Session, error) {
			sessions, err := h.Service.GetByUser(r.Context(), session.User.ID)
			if err != nil {
				return nil, err
			}
			return &sessions, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	res := make([]*responses.SessionResponse, len(*sessions))
	for i, s := range *sessions {
		res[i] = s.Response()
	}

	responses.WriteJSON(w, responses.RespondOk(res, "sessions retrieved"))
}

// TODO: finished here
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

	session, err := middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*models.Session, error) {
			dbSession, err := h.Service.GetByID(r.Context(), sid)
			if err != nil {
				return nil, err
			}

			if dbSession.User.ID != uid {
				return nil, models.ErrForbidden
			}

			return dbSession, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(session.Response(), "Session retrieved"))
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

	_, err = middlewares.RequiresOwner(
		r,
		func(requester *models.User) bool {
			return requester.Role == models.ROLE_USER && requester.ID != uid
		},
		func(session *models.Session) (*models.Session, error) {
			target, err := h.Service.GetByID(r.Context(), sid)
			if err != nil {
				return nil, err
			}

			if session.User.Role == models.ROLE_USER && target.User.ID != uid {
				return nil, models.ErrForbidden
			}

			if err := h.Service.Close(r.Context(), sid); err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "session closed"))
}

func (h *SessionHandler) deleteExpired(w http.ResponseWriter, r *http.Request) {
	_, err := middlewares.RequiresOwner(r, func(requester *models.User) bool {
		return requester.Role == models.ROLE_USER
	}, func(session *models.Session) (*any, error) {
		if err := h.Service.DeleteExpired(r.Context()); err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		handleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "expired sessions deleted"))
}
