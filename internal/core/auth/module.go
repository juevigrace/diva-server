package auth

import "github.com/juevigrace/diva-server/internal/core"

type AuthModule struct {
	handler *AuthHandler
	service *AuthService
}

func NewAuthModule(
	pService *permission.PermissionService,
	sService *session.SessionService,
	uService *user.UserService,
) core.Module {
	service := NewAuthService(pService, sService, uService)
	return &AuthModule{
		handler: NewAuthHandler(service, sService),
		service: service,
	}
}

func (m *AuthModule) Handler() core.Handler {
	return m.handler
}

func (m *AuthModule) Service() core.Service {
	return m.service
}
