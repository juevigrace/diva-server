package user

import (
	"github.com/juevigrace/diva-server/internal/core"
)

type UserModule struct {
	handler *UserHandler
	service *UserService
}

func NewUserModule(service *UserService) core.Module {
	return &UserModule{
		handler: NewUserHandler(service),
		service: service,
	}
}

func (m *UserModule) Handler() core.Handler {
	return m.handler
}

func (m *UserModule) Service() core.Service {
	return m.service
}
