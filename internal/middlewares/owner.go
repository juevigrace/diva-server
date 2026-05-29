package middlewares

import (
	"net/http"

	"github.com/juevigrace/diva-server/internal/models"
)

func RequiresOwnerOrPerms[T any](r *http.Request, predicate func(requester *models.User) bool, handler func(session *models.Session) (*T, error)) (*T, error) {
	session, ok := GetSessionFromContext(r.Context())
	if !ok {
		return nil, models.ErrSessionNotFound
	}

	if predicate(&session.User) {
		return nil, models.ErrForbidden
	}

	return handler(session)
}
