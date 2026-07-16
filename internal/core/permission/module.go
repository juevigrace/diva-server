package permission

import (
	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage"
)

type PermissionModule struct {
	Handler *PermissionHandler
	Repo    *PermissionRepo
}

func NewPermissionModule(store storage.PermissionStore) *PermissionModule {
	repo := NewPermissionRepo(store)
	return &PermissionModule{
		Handler: NewPermissionHandler(repo),
		Repo:    repo,
	}
}

func (m *PermissionModule) Routes(r chi.Router, sCall middlewares.SessionCall, uCall middlewares.UserCall) {
	r.Route("/permissions", func(p chi.Router) {
		p.Use(middlewares.RequiresSession(sCall, uCall), middlewares.RequireVerified())

		p.With(
			middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
			middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_READ),
		).Get("/", m.Handler.list)

		p.Route("/{pid}", func(pid chi.Router) {
			pid.With(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
				middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_READ),
			).Get("/", m.Handler.getByID)

			pid.With(
				middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR),
				middlewares.RequirePermission(models.PERMISSION_PERMISSIONS_WRITE),
			).Put("/", m.Handler.update)

			pid.With(
				middlewares.RequireRole(models.ROLE_ADMIN),
			).Patch("/level", m.Handler.updateRoleLevel)
		})
	})
}
