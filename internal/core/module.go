package core

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Provider[T any] interface {
	GetByID(ctx context.Context, id uuid.UUID) (T, error)
}

type Service interface {
}

type Handler interface {
	Routes(r chi.Router)
}

type Module interface {
	Handler() Handler
	Service() Service
}
