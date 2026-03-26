package di

import (
	"github.com/juevigrace/diva-server/internal/handler"
)

type HandlerModule struct {
	Auth    *handler.AuthHandler
	User    *handler.UserHandler
	Session *handler.SessionHandler
}

func NewHandlerModule(services *ServiceModule) *HandlerModule {
	actionHandler := handler.NewUserActionsHandler(services.Action)

	auth := handler.NewAuthHandler(services.Auth, services.Session)
	sessionHandler := handler.NewSessionHandler(services.Session)

	userPermission := handler.NewUserPermissionHandler(services.UserPermission)
	userPreferences := handler.NewUserPreferencesHandler(services.UserPreferences)
	userMe := handler.NewUserMeHandler(services.Action, services.User, services.Verification, userPreferences, actionHandler)
	user := handler.NewUserHandler(services.Session, services.User, userMe, userPermission)

	return &HandlerModule{
		Auth:    auth,
		User:    user,
		Session: sessionHandler,
	}
}
