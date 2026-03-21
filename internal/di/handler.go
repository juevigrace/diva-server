package di

import (
	"github.com/juevigrace/diva-server/internal/handler"
)

type HandlerModule struct {
	Auth         *handler.AuthHandler
	Verification *handler.VerificationHandler
	User         *handler.UserHandler
	Session      *handler.SessionHandler
}

func NewHandlerModule(services *ServiceModule) *HandlerModule {
	sessionHandler := handler.NewSessionHandler(services.Session)
	verification := handler.NewVerificationHandler(services.Verification)
	userPermission := handler.NewUserPermissionHandler(services.UserPermission)
	userPreferences := handler.NewUserPreferencesHandler(services.UserPreferences)
	userMe := handler.NewUserMeHandler(services.User, userPreferences)
	user := handler.NewUserHandler(services.User, services.Verification, userMe, userPermission)
	auth := handler.NewAuthHandler(services.Auth, services.Session)

	return &HandlerModule{
		Auth:         auth,
		Verification: verification,
		User:         user,
		Session:      sessionHandler,
	}
}
