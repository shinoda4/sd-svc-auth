package grpc

import (
	"log"
	"net"
	"os"

	authpb "github.com/shinoda4/sd-grpc-proto/auth/v1"
	"github.com/shinoda4/sd-svc-auth/internal/service/auth"
	"google.golang.org/grpc"
)

func RunGRPCServer(authService *auth.Service) {
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, NewAuthServer(authService))

	log.Printf("gRPC server running on %s", os.Getenv("GRPC_PORT"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
