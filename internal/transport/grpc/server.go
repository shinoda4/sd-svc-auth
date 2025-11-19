package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/shinoda4/sd-grpc-proto/auth/v1"
	"github.com/shinoda4/sd-svc-auth/internal/service/auth"
	"github.com/shinoda4/sd-svc-auth/pkg/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RunGRPCServer(authService *auth.Service) {
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			func(
				ctx context.Context,
				req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler,
			) (interface{}, error) {
				// logging interceptor
				resp, err := handler(ctx, req)
				if err != nil {
					log.Printf("[gRPC] %s error: %v", info.FullMethod, err)
				}
				return resp, err
			},
			AuthInterceptor(authService), // 认证 interceptor
		),
	)
	authpb.RegisterAuthServiceServer(grpcServer, NewAuthServer(authService))

	log.Printf("gRPC server running on %s", os.Getenv("GRPC_PORT"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func RunGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	grpcAddr := fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT"))

	// 注册 AuthService 到 grpc-gateway
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	httpAddr := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	log.Printf("HTTP gateway running on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("failed to serve HTTP gateway: %v", err)
	}
}
func AuthInterceptor(authService *auth.Service) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 白名单，不需要认证的 API
		noAuthMethods := map[string]bool{
			"/auth.v1.AuthService/Register":       true,
			"/auth.v1.AuthService/VerifyEmail":    true,
			"/auth.v1.AuthService/Login":          true,
			"/auth.v1.AuthService/ForgotPassword": true,
			"/auth.v1.AuthService/ResetPassword":  true,
		}

		if noAuthMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// 下面是需要 token 的逻辑
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		rawToken := strings.TrimPrefix(authHeaders[0], "Bearer ")
		rawToken = strings.TrimSpace(rawToken)

		claims, err := token.ParseToken(rawToken)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, "claims", claims)
		ctx = context.WithValue(ctx, "raw_token", rawToken)

		return handler(ctx, req)
	}
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthService *auth.Service
}

func NewAuthServer(authService *auth.Service) *AuthServer {
	return &AuthServer{AuthService: authService}
}
