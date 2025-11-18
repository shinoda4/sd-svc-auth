package auth

import (
	"context"
	"fmt"

	"github.com/shinoda4/sd-svc-auth/internal/service"
	"github.com/shinoda4/sd-svc-auth/internal/service/entity"
	"github.com/shinoda4/sd-svc-auth/pkg/email"
)

func (s *Service) Register(ctx context.Context, userEmail, username, password string, sendEmail bool, verifyLink string) (entity.UserEntity, string, error) {
	user, err := s.db.CreateUser(ctx, userEmail, username, password)
	if err != nil {
		return nil, "", err
	}

	verifyToken := service.GenerateVerifyToken()
	// 保存 token
	if err := s.db.SetVerifyToken(ctx, user.GetID(), verifyToken); err != nil {
		return nil, "", err
	}
	if sendEmail {
		fullLink := fmt.Sprintf("%s?token=%s", verifyLink, verifyToken)
		if err := email.SendVerifyEmail(userEmail, username, fullLink); err != nil {
			return user, "", err
		}
	}
	return user, verifyToken, nil
}
