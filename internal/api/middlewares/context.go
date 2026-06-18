package middlewares

import (
	"context"

	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/pkg/errs"
)

type contextKey string

const requestContextKey contextKey = "request_context"

type RequestContext struct {
	Session *models.Session
	Cache   map[string]any
}

func NewRequestContext(session *models.Session) *RequestContext {
	return &RequestContext{
		Session: session,
		Cache:   make(map[string]any, 0),
	}
}

func SetRequestContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, requestContextKey, rc)
}

func GetRequestContext(ctx context.Context) (*RequestContext, error) {
	rc, ok := ctx.Value(requestContextKey).(*RequestContext)
	if !ok {
		return nil, errs.ErrContextIsMissing
	}
	if rc.Session == nil {
		return nil, errs.ErrSessionNotFound
	}
	return rc, nil
}

func GetSessionFromContext(ctx context.Context) (*models.Session, bool) {
	rc, err := GetRequestContext(ctx)
	if err != nil {
		return nil, false
	}
	if rc.Session == nil {
		return nil, false
	}
	return rc.Session, true
}
