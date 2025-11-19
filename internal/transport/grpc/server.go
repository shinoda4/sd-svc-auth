package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/shinoda4/sd-grpc-proto/auth/v1"
	"github.com/shinoda4/sd-svc-auth/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunGRPCServer(authService *auth.Service) {
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("[gRPC] %s error: %v", info.FullMethod, err)
		}
		return resp, err
	}))
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
