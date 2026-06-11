package permission

import (
	"github.com/juevigrace/diva-server/internal/core"
	"github.com/juevigrace/diva-server/storage/db"
)

type PermissionModule struct {
	handler *PermissionHandler
	service *PermissionService
}

func NewPermissionModule(queries *db.Queries, sService *sesison.SessionService) core.Module {
	service := NewPermissionService(queries)
	return &PermissionModule{
		handler: NewPermissionHandler(service, sService),
		service: service,
	}
}

func (m *PermissionModule) Handler() core.Handler {
	return m.handler
}

func (m *PermissionModule) Service() core.Service {
	return m.service
}
