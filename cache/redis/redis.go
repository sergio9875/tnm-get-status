package redis

import (
	"context"
	"errors"
	"fmt"
	"malawi-getstatus/cache"
	"malawi-getstatus/models"
	"time"

	impl "github.com/go-redis/redis/v8"
)

type Holder struct {
	Cache *impl.Client
}

func NewCache(config *models.Cache) (cache.Cache, error) {
	if *config.Type != "redis" {
		//log.Errorf("ROOT", "Invalid cache type passed to redis cache service: %s", config.Type)
		return nil, errors.New("invalid cache type passed to redis cache service")
	}
	rdb := impl.NewClient(&impl.Options{
		Addr: fmt.Sprintf("%s:%d", *config.Host, *config.Port),
		DB:   *config.Database,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		//log.Error("ROOT", fmt.Sprintf("%s:%d", *config.Host, *config.Port))
		//log.Errorf("ROOT", "Redis Ping returned an error: ", err.Error())
		return nil, err
	}

	return &Holder{Cache: rdb}, nil
}

func (h *Holder) HSet(ctx context.Context, key string, field string, value string) error {
	args := make(map[string]string, 0)
	args[field] = value
	return h.Cache.HSet(ctx, key, args).Err()
}

func (h *Holder) Get(ctx context.Context, key string) (string, error) {
	return h.Cache.Get(ctx, key).Result()
}

func (h *Holder) HGet(ctx context.Context, key string, field string) (string, error) {
	return h.Cache.HGet(ctx, key, field).Result()
}

func (h *Holder) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return h.Cache.HGetAll(ctx, key).Result()
}

func (h *Holder) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return h.Cache.Expire(ctx, key, ttl).Err()
}
