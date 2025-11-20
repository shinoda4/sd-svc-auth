/*
 * Copyright (c) 2025-11-20 shinoda4
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
