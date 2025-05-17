package redis

import (
	"context"
	"time"

	"github.com/gorozcovcp/little-twitter/internal/ports"
	"github.com/redis/go-redis/v9"
)

type RedisTimelineCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisTimelineCache(client *redis.Client, ttl time.Duration) ports.TimelineCache {
	return &RedisTimelineCache{
		client: client,
		ttl:    ttl,
	}
}

func (r *RedisTimelineCache) Get(ctx context.Context, userID string) ([]byte, error) {
	val, err := r.client.Get(ctx, r.key(userID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (r *RedisTimelineCache) Set(ctx context.Context, userID string, data []byte) error {
	return r.client.Set(ctx, r.key(userID), data, r.ttl).Err()
}

func (r *RedisTimelineCache) Delete(ctx context.Context, userID string) error {
	return r.client.Del(ctx, r.key(userID)).Err()
}

func (r *RedisTimelineCache) key(userID string) string {
	return "timeline:" + userID
}
