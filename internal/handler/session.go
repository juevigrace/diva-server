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

type SessionHandler struct {
	sService *service.SessionService
}

func NewSessionHandler(sService *service.SessionService) *SessionHandler {
	return &SessionHandler{sService: sService}
}

func (h *SessionHandler) UserRoutes(r chi.Router) {
	r.Route("/sessions", func(s chi.Router) {
		s.With(middlewares.RequireResourceOwner(
			"uid",
			func(ctx context.Context, reqid, resid uuid.UUID) (any, bool) {
				return nil, reqid == resid
			},
			models.PERMISSION_SESSIONS_READ,
		)).Get("/", h.listByUser)
	})
}

func (h *SessionHandler) Routes(r chi.Router) {
	r.Route("/sessions", func(s chi.Router) {
		s.Use(middlewares.RequiresSession(h.sService.GetByID))

		s.Route("/{sid}", func(sid chi.Router) {
			sid.With(middlewares.RequireResourceOwner(
				"sid",
				func(ctx context.Context, reqid, resid uuid.UUID) (any, bool) {
					dbSession, err := h.sService.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if dbSession.User.ID != reqid {
						return nil, false
					}
					return dbSession, true
				},
				models.PERMISSION_SESSIONS_READ,
			)).Get("/", h.getByID)
			sid.With(middlewares.RequireResourceOwner(
				"sid",
				func(ctx context.Context, reqid, resid uuid.UUID) (any, bool) {
					dbSession, err := h.sService.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if dbSession.User.ID != reqid {
						return nil, false
					}
					return dbSession, true
				},
				models.PERMISSION_SESSIONS_WRITE,
			)).Delete("/", h.close)
		})

		s.Group(func(admin chi.Router) {
			admin.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
			admin.Delete("/expired", h.deleteExpired)
		})
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
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	session, ok := rc.Cache["sid"].(*models.Session)
	if !ok {
		sid, err := middlewares.GetUUIDFromURL(r, "sid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}

		session, err = h.sService.GetByID(r.Context(), sid)
		if err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondOk(session.Response(), "session retrieved"))
}

func (h *SessionHandler) close(w http.ResponseWriter, r *http.Request) {
	rc, err := middlewares.GetRequestContext(r.Context())
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	var sid uuid.UUID
	session, ok := rc.Cache["sid"].(*models.Session)
	if !ok {
		sid, err = middlewares.GetUUIDFromURL(r, "sid")
		if err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
	} else {
		sid = session.ID
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
