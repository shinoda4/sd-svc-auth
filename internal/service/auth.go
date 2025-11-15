package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/token"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, username, password string) error
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
	SetBlacklist(ctx context.Context, token string, ttl time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string) error
}

type AuthService struct {
	db    UserRepository
	cache CacheRepository
}

func NewAuthService(db UserRepository, cache CacheRepository) *AuthService {
	return &AuthService{db: db, cache: cache}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) error {
	return s.db.CreateUser(ctx, email, username, password)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, accessTTL, refreshTTL time.Duration, err error) {
	u, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", 0, 0, err
	}
	if !u.CheckPassword(password) {
		return "", "", 0, 0, ErrInvalidPassword
	}

	accessToken, accessTTL, err = token.GenerateJWT(u.GetID(), u.GetEmail())
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshToken, refreshTTL, err = token.GenerateRefreshJWT(u.GetID(), u.GetEmail())
	if err != nil {
		return "", "", 0, 0, err
	}

	// 只缓存 refresh token
	if err := s.cache.StoreToken(ctx, u.GetID(), refreshToken, refreshTTL); err != nil {
		return "", "", 0, 0, err
	}

	return accessToken, refreshToken, accessTTL, refreshTTL, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (newAccessToken string, accessTTL time.Duration, err error) {
	claims, err := token.ParseAndValidateRefresh(refreshToken)
	if err != nil {
		return "", 0, err
	}

	// 验证缓存中的 refresh token
	stored, err := s.cache.GetToken(ctx, claims.UserID)
	if err != nil || stored != refreshToken {
		return "", 0, ErrInvalidToken
	}

	// 生成新的 access token
	newAccessToken, accessTTL, err = token.GenerateJWT(claims.UserID, claims.Email)
	if err != nil {
		return "", 0, err
	}

	return newAccessToken, accessTTL, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	claims, err := token.ParseAndValidate(tokenStr)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (s *AuthService) Logout(ctx context.Context, tokenStr string) error {
	// 解析 token 类型
	claims, err := token.ParseAndValidate(tokenStr)
	if err != nil {
		log.Println("Error validating token:", err)
		return err
	}

	switch claims.TokenType {
	case "access":
		// Access token 黑名单：存储到 Redis，过期时间和 token 一样
		ttl := time.Until(claims.ExpiresAt.Time)
		return s.cache.SetBlacklist(ctx, tokenStr, ttl)

	case "refresh":
		// 删除 Redis 中的 Refresh token
		userID := claims.Subject
		return s.cache.DeleteRefreshToken(ctx, userID)

	default:
		return errors.New("unknown token type")
	}
}
