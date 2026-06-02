package middlewares

import (
	"net/http"
	"time"

	"github.com/juevigrace/diva-server/internal/models"
)

func RequiresOwner[T any](r *http.Request, predicate func(requester *models.User) bool, handler func(session *models.Session) (*T, error)) (*T, error) {
	session, ok := GetSessionFromContext(r.Context())
	if !ok {
		return nil, models.ErrSessionNotFound
	}

	if predicate(&session.User) {
		return nil, models.ErrForbidden
	}

	return handler(session)
}

func RequiresPermission(requester *models.User, permAction models.PermissionAction) error {
	if requester.Role == models.ROLE_USER {
		perm := requester.Permissions[permAction]

		exp := time.UnixMilli(perm.ExpiresAt)
		if exp.Before(time.Now().UTC()) {
			return models.ErrPermissionExpired
		}

		if !perm.Granted {
			return models.ErrPermissionDenied
		}
	}

	return nil
}
