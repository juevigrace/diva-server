package auth

import (
	"github.com/juevigrace/diva-server/internal/core"
	"github.com/juevigrace/diva-server/internal/core/permission"
	"github.com/juevigrace/diva-server/internal/core/session"
	"github.com/juevigrace/diva-server/internal/core/user"
)

type AuthModule struct {
	handler *AuthHandler
	service *AuthService
}

func NewAuthModule(
	pService *permission.PermissionService,
	sService *session.SessionService,
	uService *user.UserService,
) *AuthModule {
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
