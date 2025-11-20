/*
 * Copyright (c) 2025-11-20 shinoda4
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
