package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shinoda4/sd-svc-auth/pkg/token"
)

func (s *Service) Logout(ctx context.Context, tokenStr string) error {
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
