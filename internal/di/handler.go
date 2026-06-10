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
	session := handler.NewSessionHandler(services.Session)
	userPermission := handler.NewUserPermissionHandler(services.UserPermission)
	userPreferences := handler.NewUserPreferencesHandler(services.UserPermission, services.UserPreferences)
	userAction := handler.NewUserActionsHandler(services.UserActions, services.Session)
	userProfile := handler.NewUserProfileHandler(services.UserPermission, services.UserProfile)
	user := handler.NewUserHandler(services.User, services.UserState, session, userAction, userPermission, userPreferences, userProfile)
	verification := handler.NewVerificationHandler(services.Session, services.User, services.UserActions, services.UserPermission, services.UserState, services.Verification)
	permissions := handler.NewPermissionsHandler(services.Permission, services.Session)

	return &HandlerModule{
		Auth:         auth,
		User:         user,
		Verification: verification,
		Permissions:  permissions,
		Session:      session,
	}
}
