package grpc

import (
	"context"
	"time"

	authpb "github.com/shinoda4/sd-grpc-proto/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	rawToken, ok := ctx.Value("raw_token").(string)
	if !ok || rawToken == "" {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	newAccess, accessTTL, err := s.AuthService.Refresh(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	return &authpb.RefreshTokenResponse{
		AccessToken: newAccess,
		ExpiresIn:   timestamppb.New(time.Now().Add(accessTTL)),
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	rawToken, ok := ctx.Value("raw_token").(string)
	if !ok || rawToken == "" {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	claims, err := s.AuthService.ValidateToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}
