package di

import (
	"github.com/juevigrace/diva-server/internal/handler"
	"github.com/juevigrace/diva-server/internal/mail"
)

type HandlerModule struct {
	Verification *handler.VerificationHandler

	User    *handler.UserHandler
	Session *handler.SessionHandler
}

func NewHandlerModule(repos *RepoModule, mailClient *mail.Client) *HandlerModule {
	session := &handler.SessionHandler{Repo: repos.Session}

	verification := handler.NewVerificationHandler(repos.Verification, mailClient)

	userPermission := handler.NewUserPermissionHandler(repos.UserPermission)
	userPreferences := handler.NewUserPreferencesHandler(repos.UserPreferences)
	userMe := handler.NewUserMeHandler(repos.User, userPreferences)
	user := handler.NewUserHandler(repos.User, repos.Session, userMe, userPermission, verification)

	return &HandlerModule{
		Verification: verification,
		User:         user,
		Session:      session,
	}
}
