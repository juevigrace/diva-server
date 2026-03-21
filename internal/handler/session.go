package handler

import (
	"github.com/juevigrace/diva-server/internal/service"
)

type SessionHandler struct {
	Service *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{Service: svc}
}
