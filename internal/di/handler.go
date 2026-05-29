package di

import (
	"github.com/juevigrace/diva-server/internal/handler"
)

type HandlerModule struct {
	Auth         *handler.AuthHandler
	User         *handler.UserHandler
	Verification *handler.UserVerificationHandler
	Permissions  *handler.PermissionsHandler
}

func NewHandlerModule(services *ServiceModule) *HandlerModule {
	auth := handler.NewAuthHandler(services.Auth, services.Session)
	session := handler.NewSessionHandler(services.Session)
	userPermission := handler.NewUserPermissionHandler(services.UserPermission)
	userPreferences := handler.NewUserPreferencesHandler(services.UserPreferences)
	userAction := handler.NewUserActionsHandler(services.UserActions, services.Session)
	userProfile := handler.NewUserProfileHandler(services.UserProfile)
	user := handler.NewUserHandler(services.Session, services.User, session, userAction, userPermission, userPreferences, userProfile)
	verification := handler.NewVerificationHandler(services.Session, services.User, services.Verification)
	permissions := handler.NewPermissionsHandler(services.Permission)

	return &HandlerModule{
		Auth:         auth,
		User:         user,
		Verification: verification,
		Permissions:  permissions,
	}
}
