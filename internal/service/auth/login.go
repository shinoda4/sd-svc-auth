package auth

import (
	"context"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/service"
	"github.com/shinoda4/sd-svc-auth/pkg/token"
)

func (s *Service) Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, accessTTL, refreshTTL time.Duration, err error) {
	u, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", 0, 0, err
	}
	if !u.CheckPassword(password) {
		return "", "", 0, 0, service.ErrInvalidPassword
	}
	if !u.GetEmailVerified() {
		return "", "", 0, 0, service.ErrEmailNotVerified
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
