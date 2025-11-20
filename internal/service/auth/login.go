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
