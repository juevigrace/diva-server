package di

import (
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/service"
)

type ServiceModule struct {
	User            *service.UserService
	Session         *service.SessionService
	Auth            *service.AuthService
	Verification    *service.UserVerificationService
	UserPermission  *service.UserPermissionService
	UserPreferences *service.UserPreferencesService
	UserActions     *service.UserActionsService
	UserProfile     *service.UserProfileService
	Permission      *service.PermissionService
}

func NewServiceModule(repos *RepoModule, mailClient *mail.Client) *ServiceModule {
	session := service.NewSessionService(repos.Session)
	uAction := service.NewUserActionsService(repos.Action)
	user := service.NewUserService(repos.User, uAction)
	userPermission := service.NewUserPermissionService(repos.UserPermission)
	userPreferences := service.NewUserPreferencesService(repos.UserPreferences)
	userProfile := service.NewUserProfileService(repos.UserProfile)
	permission := service.NewPermissionService(repos.Permissions)
	verification := service.NewVerificationService(mailClient, repos.Verification, session, uAction)
	auth := service.NewAuthService(user, verification, session)

	return &ServiceModule{
		User:            user,
		Session:         session,
		Auth:            auth,
		Verification:    verification,
		UserPermission:  userPermission,
		UserPreferences: userPreferences,
		UserActions:     uAction,
		UserProfile:     userProfile,
		Permission:      permission,
	}
}
