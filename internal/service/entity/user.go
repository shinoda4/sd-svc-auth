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
}

type UserEntity interface {
	GetID() string
	GetEmail() string
	GetUsername() string
	GetEmailVerified() bool
	CheckPassword(password string) bool
}

type CacheRepository interface {
	StoreToken(ctx context.Context, userID, token string, ttl time.Duration) error
	GetToken(ctx context.Context, userID string) (string, error)
	SetBlacklist(ctx context.Context, token string, ttl time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string) error
}
