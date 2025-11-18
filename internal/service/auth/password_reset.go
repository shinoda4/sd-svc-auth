package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/service"
	"github.com/shinoda4/sd-svc-auth/pkg/email"
)

func (s *Service) PasswordReset(ctx context.Context, email_addr string, username string) error {

	user, err := s.db.GetUserByEmail(ctx, email_addr)
	if err != nil {
		return err
	}

	if username != user.GetUsername() {
		return service.ErrUsernameNotValid
	}

	emailAddress := os.Getenv("EMAIL_ADDRESS")
	if emailAddress == "" {
		return errors.New("EMAIL_ADDRESS environment variable not set")
	}

	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
	resetToken := hex.EncodeToString(tokenBytes)
	expire := time.Now().Add(time.Hour * 1)

	if err := s.db.SaveResetToken(ctx, user.GetID(), resetToken, expire); err != nil {
		return err
	}

	baseURL := os.Getenv("RESET_PASSWORD_URL")
	if baseURL == "" {
		return errors.New("RESET_PASSWORD_URL environment variable not set")
	}
	fullLink := fmt.Sprintf("%s?token=%s", baseURL, resetToken)

	body := fmt.Sprintf(
		"Dear <b>%s</b>,<br><br>Please click the following link to reset your password:<br><a href='%s'>Reset Password</a><br><br>If you did not request this, please ignore this email.",
		username, fullLink,
	)

	return email.SendEmail(emailAddress, email_addr, "Reset your password!", body)
}

func (s *Service) PasswordResetConfirm(ctx context.Context, token, newPassword string) error {
	user, err := s.db.GetUserByResetToken(ctx, token)
	if err != nil {
		return err
	}

	if time.Now().After(user.GetResetTokenExpire()) {
		return errors.New("Token expired!")
	}

	err = s.db.UpdatePassword(ctx, user.GetID(), newPassword)
	if err != nil {
		return err
	}

	return s.db.ClearResetToken(ctx, user.GetID())
}
