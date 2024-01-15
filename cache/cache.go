package cache

import (
	"context"
	"time"
)

type Cache interface {
	HSet(context.Context, string, string, string) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Expire(context.Context, string, time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
