package user

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/api/core/permission"
	"github.com/juevigrace/diva-server/internal/api/core/session"
	"github.com/juevigrace/diva-server/internal/api/core/user/actions"
	"github.com/juevigrace/diva-server/internal/api/core/user/permissions"
	"github.com/juevigrace/diva-server/internal/api/core/user/preferences"
	"github.com/juevigrace/diva-server/internal/api/core/user/profile"
	"github.com/juevigrace/diva-server/internal/api/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/pkg/filehelper"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserModule struct {
	sRepo *session.SessionRepo
	sHandler *session.SessionHandler

	uHandler    *UserHandler
	uaHandler   *actions.UserActionsHandler
	upHandler   *permissions.UserPermissionHandler
	uprHandler  *preferences.UserPreferencesHandler
	uproHandler *profile.UserProfileHandler

	URepo  *UserRepo
	UARepo *actions.UserActionsRepo
	UPRepo *permissions.UserPermissionRepo
	USRepo *UserStateRepo
}

func NewUserModule(
	queries *db.Queries,
	pRepo *permission.PermissionRepo,
	sRepo *session.SessionRepo,
	sHandler *session.SessionHandler,
	files *filehelper.FileHelper,
) *UserModule {
	uaRepo := actions.NewUserActionsRepo(queries)
	upRepo := permissions.NewUserPermissionRepo(queries, pRepo)
	uprRepo := preferences.NewUserPreferencesRepo(queries, upRepo)
	uproRepo := profile.NewUserProfileRepo(queries, upRepo)

	uaHandler := actions.NewUserActionsHandler(uaRepo)
	upHandler := permissions.NewUserPermissionHandler(upRepo)
	uprHandler := preferences.NewUserPreferencesHandler(upRepo, uprRepo)
	uproHandler := profile.NewUserProfileHandler(uproRepo, files)

	usRepo := NewUserStateRepo(queries)
	uRepo := NewUserRepo(queries, sRepo, uaRepo, upRepo, uproRepo, usRepo)
	uHandler := NewUserHandler(uRepo, usRepo)

	return &UserModule{
		sRepo:    sRepo,
		sHandler:    sHandler,
		uHandler:    uHandler,
		uaHandler:   uaHandler,
		upHandler:   upHandler,
		uprHandler:  uprHandler,
		uproHandler: uproHandler,
		URepo:    uRepo,
		UARepo:   uaRepo,
		UPRepo:   upRepo,
		USRepo:   usRepo,
	}
}

func (m *UserModule) Routes(r chi.Router) {
	r.Route("/user", func(u chi.Router) {
		u.Route("/check", func(check chi.Router) {
			check.Get("/username/{username}", m.uHandler.checkUsername)
			check.Get("/email/{email}", m.uHandler.checkEmail)
		})

		u.Group(func(auth chi.Router) {
			auth.Use(middlewares.RequiresSession(m.sRepo.GetByID, m.URepo.GetByID))

			auth.Group(func(admin chi.Router) {
				admin.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
				admin.Get("/", m.uHandler.getAll)
				admin.Post("/", m.uHandler.create)
			})

			auth.Route("/{uid}", func(uid chi.Router) {
				uid.Get("/", m.uHandler.getByID)

				uid.With(
					middlewares.RequirePermission(models.PERMISSION_USERS_EMAIL_WRITE),
					middlewares.RequireResourceOwner(
						&middlewares.RequireOwnerParams{
							UrlParams: []string{"uid"},
							Perms:     []models.PermissionAction{models.PERMISSION_USERS_EMAIL_WRITE},
						},
						func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
							resid, err := uuid.Parse(resParams[0])
							if err != nil {
								return nil, false
							}
							return map[string]any{"uid": resid}, reqid == resid
						},
					),
				).Patch("/email", m.uHandler.updateEmail)
				uid.With(
					middlewares.RequirePermission(models.PERMISSION_USERS_PHONE_WRITE),
					middlewares.RequireResourceOwner(
						&middlewares.RequireOwnerParams{
							UrlParams: []string{"uid"},
							Perms:     []models.PermissionAction{models.PERMISSION_USERS_PHONE_WRITE},
						},
						func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
							resid, err := uuid.Parse(resParams[0])
							if err != nil {
								return nil, false
							}
							return map[string]any{"uid": resid}, reqid == resid
						},
					),
				).Patch("/phone", m.uHandler.updatePhone)
				uid.With(
					middlewares.RequirePermission(models.PERMISSION_USERS_USERNAME_WRITE),
					middlewares.RequireResourceOwner(
						&middlewares.RequireOwnerParams{
							UrlParams: []string{"uid"},
							Perms:     []models.PermissionAction{models.PERMISSION_USERS_USERNAME_WRITE},
						},
						func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
							resid, err := uuid.Parse(resParams[0])
							if err != nil {
								return nil, false
							}
							return map[string]any{"uid": resid}, reqid == resid
						},
					),
				).Patch("/username", m.uHandler.updateUsername)
				uid.With(
					middlewares.RequirePermission(models.PERMISSION_USERS_PASSWORD_WRITE),
					middlewares.RequireResourceOwner(
						&middlewares.RequireOwnerParams{
							UrlParams: []string{"uid"},
							Perms:     []models.PermissionAction{models.PERMISSION_USERS_PASSWORD_WRITE},
						},
						func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
							resid, err := uuid.Parse(resParams[0])
							if err != nil {
								return nil, false
							}
							return map[string]any{"uid": resid}, reqid == resid
						},
					),
				).Patch("/password", m.uHandler.updatePassword)

				uid.Group(func(admin chi.Router) {
					admin.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
					admin.With(middlewares.RequirePermission(models.PERMISSION_USERS_ROLE_WRITE)).Patch("/role", m.uHandler.updateRole)
					admin.With(middlewares.RequirePermission(models.PERMISSION_USERS_RESTORE_WRITE)).Patch("/restore", m.uHandler.restore)
				})

				uid.Group(func(wg chi.Router) {
					wg.Use(middlewares.RequireResourceOwner(
						&middlewares.RequireOwnerParams{
							UrlParams: []string{"uid"},
							Perms:     []models.PermissionAction{models.PERMISSION_USERS_WRITE},
						},
						func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
							resid, err := uuid.Parse(resParams[0])
							if err != nil {
								return nil, false
							}
							return map[string]any{"uid": resid}, reqid == resid
						},
					))
					wg.Delete("/", m.uHandler.softDelete)
					wg.Delete("/forever", m.uHandler.delete)
				})

				uid.Route("/status", func(sr chi.Router) {
					sr.With(middlewares.RequireResourceOwner(
						&middlewares.RequireOwnerParams{
							UrlParams: []string{"uid"},
						},
						func(_ context.Context, reqid uuid.UUID, resParams []string) (map[string]any, bool) {
							resid, err := uuid.Parse(resParams[0])
							if err != nil {
								return nil, false
							}
							return map[string]any{"uid": resid}, reqid == resid
						},
					)).Post("/ping", m.uHandler.pingStatus)

					sr.Group(func(admin chi.Router) {
						admin.Use(middlewares.RequireRole(models.ROLE_ADMIN, models.ROLE_MODERATOR))
						admin.With(middlewares.RequirePermission(models.PERMISSION_USERS_VERIFIED_WRITE)).Patch("/verified", m.uHandler.updateVerified)
						admin.With(middlewares.RequirePermission(models.PERMISSION_USERS_WRITE)).Put("/", m.uHandler.updateStatus)
					})
				})

				m.uaHandler.UserRoutes(uid)
				m.upHandler.UserRoutes(uid)
				m.uprHandler.UserRoutes(uid)
				m.uproHandler.UserRoutes(uid)
				m.sHandler.UserRoutes(uid)
			})

			m.uaHandler.Routes(auth)
			m.uprHandler.Routes(auth)
		})
	})
}
