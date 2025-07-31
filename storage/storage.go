package storage

import "context"

type Storage interface {
	Increment(ctx context.Context, key string) (int64, error)

	SetExpiration(ctx context.Context, key string, duration int) error

	GetCounter(ctx context.Context, key string) (int64, error)

	IsBlocked(ctx context.Context, key string) (bool, error)
}
