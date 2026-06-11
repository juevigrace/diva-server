package user

import (
	"github.com/juevigrace/diva-server/internal/core"
	"github.com/juevigrace/diva-server/storage/db"
)

type UserModule struct {
	handler *UserHandler
	service *UserService
}

func NewUserModule(
	queries *db.Queries,
) core.Module {
	service := NewUserService(queries)
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
