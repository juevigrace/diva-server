package session

import (
	"github.com/juevigrace/diva-server/storage/db"
)

type SessionModule struct {
	handler *SessionHandler
	service *SessionService
}

func NewSessionModule(
	queries *db.Queries,
) *SessionModule {
	service := NewSessionService(queries)
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
