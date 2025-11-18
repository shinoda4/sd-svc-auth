package auth

import (
	"context"
	"errors"
	"fmt"
	"os"

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

	if err := s.db.SetEmailVerified(ctx, user.GetID()); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	if sendEmail {
		subject := fmt.Sprintf("Welcome! %s", user.GetUsername())
		body := "Dear <b>" + user.GetUsername() + "</b>, you are already verified! Welcome to our system!"
		emailAddress := os.Getenv("EMAIL_ADDRESS")
		if emailAddress == "" {
			return errors.New("EMAIL_ADDRESS environment variable not set")
		}

		welcomeErr := email.SendEmail(emailAddress, user.GetEmail(), subject, body)
		if welcomeErr != nil {
			return welcomeErr
		}
	}

	return nil
}
