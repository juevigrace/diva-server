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
	userPreferences := handler.NewUserPreferencesHandler(services.UserPreferences, services.UserPermission)
	userAction := handler.NewUserActionsHandler(services.UserActions, services.Session)
	userProfile := handler.NewUserProfileHandler(services.UserProfile, services.UserPermission)
	user := handler.NewUserHandler(services.User, session, userAction, userPermission, userPreferences, userProfile)
	verification := handler.NewVerificationHandler(services.Session, services.User, services.UserActions, services.Verification, services.UserPermission)
	permissions := handler.NewPermissionsHandler(services.Permission, services.Session)

	return &HandlerModule{
		Auth:         auth,
		User:         user,
		Verification: verification,
		Permissions:  permissions,
		Session:      session,
	}
}
