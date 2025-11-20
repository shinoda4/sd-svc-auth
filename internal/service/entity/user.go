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

package entity

import (
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, username, password string) (UserEntity, error)
	GetUserByEmail(ctx context.Context, email string) (UserEntity, error)
	SetVerifyToken(ctx context.Context, userID, token string) error
	GetUserByVerifyToken(ctx context.Context, token string) (UserEntity, error)
	SetEmailVerified(ctx context.Context, userID string) error
	SaveResetToken(ctx context.Context, s string, resetToken string, expire time.Time) error
	GetUserByResetToken(ctx context.Context, token string) (UserEntity, error)
	UpdatePassword(ctx context.Context, userID, newPassword string) error
	ClearResetToken(ctx context.Context, userID string) error
}

type UserEntity interface {
	GetID() string
	GetEmail() string
	GetUsername() string
	GetEmailVerified() bool
	CheckPassword(password string) bool
	GetResetTokenExpire() time.Time
}

type CacheRepository interface {
	StoreToken(ctx context.Context, userID, token string, ttl time.Duration) error
	GetToken(ctx context.Context, userID string) (string, error)
	SetBlacklist(ctx context.Context, token string, ttl time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string) error
}
