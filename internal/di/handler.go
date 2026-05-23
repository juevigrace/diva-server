package di

import (
	"github.com/juevigrace/diva-server/internal/handler"
)

type HandlerModule struct {
	Auth         *handler.AuthHandler
	User         *handler.UserHandler
	Session      *handler.SessionHandler
	Verification *handler.UserVerificationHandler
	Permissions  *handler.PermissionsHandler
}

func NewHandlerModule(services *ServiceModule) *HandlerModule {
	auth := handler.NewAuthHandler(services.Auth, services.Session)
	sessionHandler := handler.NewSessionHandler(services.Session)

	userPermission := handler.NewUserPermissionHandler(services.UserPermission, services.Session)
	userPreferences := handler.NewUserPreferencesHandler(services.UserPreferences, services.Session)
	uActionHandler := handler.NewUserActionsHandler(services.UserActions, services.Session)
	userMe := handler.NewUserMeHandler(services.User, services.UserProfile, userPreferences, uActionHandler)
	user := handler.NewUserHandler(services.Session, services.User, userMe, userPermission)

	verification := handler.NewVerificationHandler(services.Session, services.User, services.Verification)
	permissions := handler.NewPermissionsHandler(services.Permission, services.Session)

	return &HandlerModule{
		Auth:         auth,
		User:         user,
		Session:      sessionHandler,
		Verification: verification,
		Permissions:  permissions,
	}
}
