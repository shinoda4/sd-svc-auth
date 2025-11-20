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
