package di

import (
	"github.com/juevigrace/diva-server/internal/mail"
	"github.com/juevigrace/diva-server/internal/service"
	"github.com/juevigrace/diva-server/storage/db"
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

func NewServiceModule(queries *db.Queries, mailClient *mail.Client) *ServiceModule {
	uAction := service.NewUserActionsService(queries)
	userPreferences := service.NewUserPreferencesService(queries)
	userProfile := service.NewUserProfileService(queries)
	permission := service.NewPermissionService(queries)
	userPermission := service.NewUserPermissionService(queries, permission)
	verification := service.NewVerificationService(mailClient, queries, uAction)
	user := service.NewUserService(queries, uAction, userPermission, userProfile, verification)
	session := service.NewSessionService(queries, user)
	auth := service.NewAuthService(user, uAction, verification, session)

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
