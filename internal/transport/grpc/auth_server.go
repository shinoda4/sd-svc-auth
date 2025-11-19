package grpc

import (
	"context"
	"fmt"
	"os"
	"time"

	authpb "github.com/shinoda4/sd-grpc-proto/auth/v1"
	"github.com/shinoda4/sd-svc-auth/internal/service/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthService *auth.Service
}

func NewAuthServer(authService *auth.Service) *AuthServer {
	return &AuthServer{AuthService: authService}
}

// Login 示例
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

// Register 示例
func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	baseURL := os.Getenv("SERVER_HOST")
	verifyLink := fmt.Sprintf("%s/api/v1/verify", baseURL)

	// 调用业务逻辑
	user, verifyToken, err := s.AuthService.Register(ctx, req.Email, req.Username, req.Password, true, verifyLink)
	if err != nil {
		return nil, err
	}

	// 返回 gRPC 响应
	return &authpb.RegisterResponse{
		UserId:      user.GetID(),
		Message:     "registered",
		VerifyToken: verifyToken,
	}, nil
}
