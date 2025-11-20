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
	"fmt"
	"os"
	"time"

	authpb "github.com/shinoda4/sd-grpc-proto/proto/auth/v1"
	"github.com/shinoda4/sd-svc-auth/pkg/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	port := os.Getenv("SERVER_PORT")
	verifyLink := fmt.Sprintf("%s/api/v1/verify", baseURL+":"+port)

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

func (s *AuthServer) VerifyEmail(ctx context.Context, req *authpb.VerifyEmailRequest) (*authpb.VerifyEmailResponse, error) {
	verifyToken := req.Token

	sendEmailBool := req.SendEmail

	if verifyToken == "" {
		return nil, status.Error(codes.InvalidArgument, "verify token is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.AuthService.VerifyEmail(ctx, verifyToken, !sendEmailBool); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "email verification failed: %v", err)
	}

	return &authpb.VerifyEmailResponse{
		Message: "email verified",
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	rawToken, ok := ctx.Value("raw_token").(string)
	if !ok || rawToken == "" {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	if err := s.AuthService.Logout(ctx, rawToken); err != nil {
		return nil, err
	}

	return &authpb.LogoutResponse{
		Message: "logout successful",
	}, nil
}

func (s *AuthServer) Me(ctx context.Context, req *authpb.MeRequest) (*authpb.MeResponse, error) {
	claims, ok := ctx.Value("claims").(*token.Claims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	return &authpb.MeResponse{
		UserId:    claims.UserID,
		Email:     claims.Email,
		ExpiresIn: timestamppb.New(time.Unix(claims.ExpiresAt.Time.Unix(), 0)),
		IssuedAt:  timestamppb.New(time.Unix(claims.IssuedAt.Time.Unix(), 0)),
	}, nil
}

func (s *AuthServer) ForgotPassword(ctx context.Context, req *authpb.ForgotPasswordRequest) (*authpb.ForgotPasswordResponse, error) {
	err := s.AuthService.PasswordReset(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	return &authpb.ForgotPasswordResponse{
		Message: "reset password email sent",
	}, nil

}

func (s *AuthServer) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.ResetPasswordResponse, error) {

	t := req.Token

	if t == "" {
		return nil, status.Error(codes.InvalidArgument, "token not provided")
	}

	if req.NewPassword != req.NewPasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "new password confirmation does not match")
	}

	err := s.AuthService.PasswordResetConfirm(ctx, t, req.NewPasswordConfirm)
	if err != nil {
		return nil, err
	}
	return &authpb.ResetPasswordResponse{
		Message: "password reset done!",
	}, nil

}
