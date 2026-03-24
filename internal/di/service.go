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
	Action          *service.ActionService
}

func NewServiceModule(repos *RepoModule, mailClient *mail.Client) *ServiceModule {
	session := service.NewSessionService(repos.Session)
	verification := service.NewVerificationService(repos.Verification, mailClient)
	user := service.NewUserService(repos.User, verification)
	userPermission := service.NewUserPermissionService(repos.UserPermission)
	userPreferences := service.NewUserPreferencesService(repos.UserPreferences)
	action := service.NewActionService(repos.Action)
	auth := service.NewAuthService(user, session, verification, action)

	return &ServiceModule{
		User:            user,
		Session:         session,
		Auth:            auth,
		Verification:    verification,
		UserPermission:  userPermission,
		UserPreferences: userPreferences,
		Action:          action,
	}
}
