package auth

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/shinoda4/sd-svc-auth/internal/service/entity"
	"github.com/shinoda4/sd-svc-auth/pkg/email"
	"github.com/shinoda4/sd-svc-auth/pkg/token"
)

func (s *Service) Register(ctx context.Context, userEmail, username, password string, sendEmail bool, verifyLink string) (entity.UserEntity, string, error) {
	user, err := s.db.CreateUser(ctx, userEmail, username, password)
	if err != nil {
		return nil, "", err
	}

	verifyToken := token.GenerateVerifyToken()
	if err := s.db.SetVerifyToken(ctx, user.GetID(), verifyToken); err != nil {
		return nil, "", err
	}
	if sendEmail {
		subject := "Verify your email!"
		fullLink := fmt.Sprintf("%s?token=%s", verifyLink, verifyToken)
		body := fmt.Sprintf("Dear <b>%s</b>, please finish your account validation by clicking the following link: <a href='%s'>Verify Email</a>", user.GetUsername(), fullLink)

		emailAddress := os.Getenv("EMAIL_ADDRESS")
		if emailAddress == "" {
			return user, "", errors.New("EMAIL_ADDRESS environment variable not set")
		}

		err := email.SendEmail(emailAddress, user.GetEmail(), subject, body)
		if err != nil {
			return user, "", err
		}
	}
	return user, verifyToken, nil
}
