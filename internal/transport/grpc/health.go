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
