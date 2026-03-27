package di

import (
	"github.com/juevigrace/diva-server/internal/handler"
)

type HandlerModule struct {
	Auth         *handler.AuthHandler
	User         *handler.UserHandler
	Session      *handler.SessionHandler
	Verification *handler.VerificationHandler
}

func NewHandlerModule(services *ServiceModule) *HandlerModule {
	auth := handler.NewAuthHandler(services.Auth, services.Session)
	sessionHandler := handler.NewSessionHandler(services.Session)

	userPermission := handler.NewUserPermissionHandler(services.UserPermission)
	userPreferences := handler.NewUserPreferencesHandler(services.UserPreferences)
	uActionHandler := handler.NewUserActionsHandler(services.UserActions)
	userMe := handler.NewUserMeHandler(services.User, userPreferences, uActionHandler)
	user := handler.NewUserHandler(services.Session, services.User, userMe, userPermission)

	verification := handler.NewVerificationHandler(services.Session, services.Verification)

	return &HandlerModule{
		Auth:         auth,
		User:         user,
		Session:      sessionHandler,
		Verification: verification,
	}
}
