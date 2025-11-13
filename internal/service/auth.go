package service

import (
	"context"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/token"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, password string) error
	GetUserByEmail(ctx context.Context, email string) (UserEntity, error)
}

type UserEntity interface {
	GetID() string
	GetEmail() string
	CheckPassword(password string) bool
}

type CacheRepository interface {
	StoreToken(ctx context.Context, userID, token string, ttl time.Duration) error
	GetToken(ctx context.Context, userID string) (string, error)
}

type AuthService struct {
	db    UserRepository
	cache CacheRepository
}

func NewAuthService(db UserRepository, cache CacheRepository) *AuthService {
	return &AuthService{db: db, cache: cache}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	return s.db.CreateUser(ctx, email, password)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, time.Duration, error) {
	u, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", 0, err
	}
	if !u.CheckPassword(password) {
		return "", 0, err
	}
	tok, ttl, err := token.GenerateJWT(u.GetID(), u.GetEmail())
	if err != nil {
		return "", 0, err
	}
	// store token in redis with ttl
	_ = s.cache.StoreToken(ctx, u.GetID(), tok, ttl)
	return tok, ttl, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	claims, err := token.ParseAndValidate(tokenStr)
	if err != nil {
		return nil, err
	}
	stored, err := s.cache.GetToken(ctx, claims.UserID)
	if err == nil && stored != tokenStr {
		return nil, err
	}
	return claims, nil
}
