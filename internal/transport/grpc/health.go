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

package grpc

import (
	"context"

	authpb "github.com/shinoda4/sd-grpc-proto/proto/auth/v1"
)

func (s *AuthServer) HealthCheck(ctx context.Context, req *authpb.HealthCheckRequest) (*authpb.HealthCheckResponse, error) {
	return &authpb.HealthCheckResponse{
		Status: "ok",
	}, nil
}
