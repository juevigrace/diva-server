package di

import (
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/service"
)

type ServiceModule struct {
	User            *service.UserService
	Session         *service.SessionService
	Auth            *service.AuthService
	Verification    *service.VerificationService
	UserPermission  *service.UserPermissionService
	UserPreferences *service.UserPreferencesService
	UserActions     *service.UserActionsService
}

func NewServiceModule(repos *RepoModule, mailClient *mail.Client) *ServiceModule {
	session := service.NewSessionService(repos.Session)
	uAction := service.NewUserActionsService(repos.Action)
	user := service.NewUserService(repos.User, uAction)
	userPermission := service.NewUserPermissionService(repos.UserPermission)
	userPreferences := service.NewUserPreferencesService(repos.UserPreferences)
	verification := service.NewVerificationService(mailClient, repos.Verification, session, user, uAction)
	auth := service.NewAuthService(user, session)

	return &ServiceModule{
		User:            user,
		Session:         session,
		Auth:            auth,
		Verification:    verification,
		UserPermission:  userPermission,
		UserPreferences: userPreferences,
		UserActions:     uAction,
	}
}
