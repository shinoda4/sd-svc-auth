package repo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) SetBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := "blacklist:" + token
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *RedisCache) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := "refresh_token:" + userID
	return r.client.Del(ctx, key).Err()
}

func NewRedis(addr, password string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &RedisCache{client: rdb}
}

func (r *RedisCache) StoreToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	return r.client.Set(ctx, "token:"+userID, token, ttl).Err()
}

func (r *RedisCache) GetToken(ctx context.Context, userID string) (string, error) {
	return r.client.Get(ctx, "token:"+userID).Result()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
