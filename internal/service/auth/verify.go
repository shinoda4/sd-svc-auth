package auth

import (
	"context"
	"fmt"

	"github.com/shinoda4/sd-svc-auth/pkg/email"
	"github.com/shinoda4/sd-svc-auth/pkg/token"
)

func (s *Service) ValidateToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	claims, err := token.ParseAndValidate(tokenStr)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string, sendEmail bool) error {
	user, err := s.db.GetUserByVerifyToken(ctx, token)
	if err != nil {
		return err
	}

	// 设置邮箱已验证
	if err := s.db.SetEmailVerified(ctx, user.GetID()); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	if sendEmail {
		welcomeErr := email.SendWelcomeEmail(user.GetEmail(), user.GetUsername())
		if welcomeErr != nil {
			return welcomeErr
		}
	}

	return nil
}
