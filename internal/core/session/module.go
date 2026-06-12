package session

import (
	"github.com/juevigrace/diva-server/internal/core"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionModule struct {
	handler *SessionHandler
	service *SessionService
}

func NewSessionModule(
	queries *db.Queries,
	provider core.Provider[*models.User],
) *SessionModule {
	service := NewSessionService(queries, provider)
	return &SessionModule{
		handler: NewSessionHandler(service),
		service: service,
	}
}

func (m *SessionModule) Handler() *SessionHandler {
	return m.handler
}

func (m *SessionModule) Service() *SessionService {
	return m.service
}
