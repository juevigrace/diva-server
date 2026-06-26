package permission

import (
	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage"
)

type PermissionModule struct {
	Handler *PermissionHandler
	Repo *PermissionRepo
}

func NewPermissionModule(store storage.PermissionStore) *PermissionModule {
	repo := NewPermissionRepo(store)
	return &PermissionModule{
		Handler: NewPermissionHandler(repo),
		Repo: repo,
	}
}

func (m *PermissionModule) Routes(r chi.Router, sCall middlewares.SessionCall, uCall middlewares.UserCall) {
	r.Route("/permissions", func(p chi.Router) {
		p.Use(middlewares.RequiresSession(sCall, uCall))

		p.With(
			middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_READ),
		).Get("/", m.Handler.list)

		p.Route("/{pid}", func(pid chi.Router) {
			pid.With(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
				middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_READ),
			).Get("/", m.Handler.getByID)

			pid.Group(func(wg chi.Router) {
				wg.Use(
					middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
					middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_WRITE),
				)
				wg.Put("/", m.Handler.update)
				wg.Patch("/restore", m.Handler.restore)
			})
			pid.Group(func(wg chi.Router) {
				wg.Use(middlewares.RequireRole(models.ROLE_ADMIN))
				wg.Delete("/", m.Handler.softDelete)
				wg.Delete("/forever", m.Handler.delete)
			})
		})

		p.With(
			middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_WRITE),
		).Post("/", m.Handler.create)
	})
}
