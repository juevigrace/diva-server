package session

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/api/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionModule struct {
	Handler *SessionHandler
	Repo *SessionRepo
}

func NewSessionModule(queries *db.Queries) *SessionModule {
	repo := NewSessionRepo(queries)
	return &SessionModule{
		Handler: NewSessionHandler(repo),
		Repo: repo,
	}
}

func (m *SessionModule) Routes(r chi.Router, uCall middlewares.UserCall) {
	r.Route("/sessions", func(s chi.Router) {
		s.Use(middlewares.RequiresSession(m.Repo.GetByID, uCall))

		s.Route("/{sid}", func(sid chi.Router) {
			sid.With(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"sid"},
					Perms:     []models.PermissionAction{models.PERMISSION_SESSIONS_READ},
				},
				func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					dbSession, err := m.Repo.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if dbSession.User.ID != reqid {
						return nil, false
					}
					return map[string]any{"sid": dbSession}, true
				},
			)).Get("/", m.Handler.getByID)
			sid.With(middlewares.RequireResourceOwner(
				&middlewares.RequireOwnerParams{
					UrlParams: []string{"sid"},
					Perms:     []models.PermissionAction{models.PERMISSION_SESSIONS_WRITE},
				},
				func(ctx context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
					resid, err := uuid.Parse(resParams[0])
					if err != nil {
						return nil, false
					}
					dbSession, err := m.Repo.GetByID(ctx, resid)
					if err != nil {
						return nil, false
					}
					if dbSession.User.ID != reqid {
						return nil, false
					}
					return map[string]any{"sid": dbSession}, true
				},
			)).Delete("/", m.Handler.close)
		})

		s.Group(func(admin chi.Router) {
			admin.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
			admin.Delete("/expired", m.Handler.deleteExpired)
		})
	})
}
