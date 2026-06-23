package storage

import (
	"context"
)

type Storage[T any] interface {
	Close() error
	Health(ctx context.Context) HealthResult
	Queries() *T
}
