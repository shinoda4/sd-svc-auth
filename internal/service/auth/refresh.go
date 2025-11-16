package auth

import (
	"context"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/service"
	"github.com/shinoda4/sd-svc-auth/pkg/token"
)

func (s *Service) Refresh(ctx context.Context, refreshToken string) (newAccessToken string, accessTTL time.Duration, err error) {
	claims, err := token.ParseAndValidateRefresh(refreshToken)
	if err != nil {
		return "", 0, err
	}

	// 验证缓存中的 refresh token
	stored, err := s.cache.GetToken(ctx, claims.UserID)
	if err != nil || stored != refreshToken {
		return "", 0, service.ErrInvalidToken
	}

	// 生成新的 access token
	newAccessToken, accessTTL, err = token.GenerateJWT(claims.UserID, claims.Email)
	if err != nil {
		return "", 0, err
	}

	return newAccessToken, accessTTL, nil
}
