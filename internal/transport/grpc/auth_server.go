package grpc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	authpb "github.com/shinoda4/sd-grpc-proto/auth/v1"
	"github.com/shinoda4/sd-svc-auth/internal/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthService *auth.Service
}

func NewAuthServer(authService *auth.Service) *AuthServer {
	return &AuthServer{AuthService: authService}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	accessToken, refreshToken, accessTTL, refreshTTL, err := s.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &authpb.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        timestamppb.New(time.Now().Add(accessTTL)),
		RefreshExpiresIn: timestamppb.New(time.Now().Add(refreshTTL)),
	}, nil
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	baseURL := os.Getenv("SERVER_HOST")
	verifyLink := fmt.Sprintf("%s/api/v1/verify", baseURL)

	user, verifyToken, err := s.AuthService.Register(ctx, req.Email, req.Username, req.Password, true, verifyLink)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{
		UserId:      user.GetID(),
		Message:     "registered",
		VerifyToken: verifyToken,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	token := strings.TrimPrefix(authHeaders[0], "Bearer ")
	token = strings.TrimSpace(token)
	if err := s.AuthService.Logout(ctx, token); err != nil {
		return nil, err
	}

	return &authpb.LogoutResponse{
		Message: "logout successful",
	}, nil
}
