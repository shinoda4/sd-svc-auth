package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/handler"
	"github.com/shinoda4/sd-svc-auth/internal/repo"
	"github.com/shinoda4/sd-svc-auth/internal/service"
	"github.com/shinoda4/sd-svc-auth/pkg/logger"
)

func main() {
	logger.Init()

	// load env or provide defaults
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://sd_auth:sd_pass@postgres:5432/sd_auth?sslmode=disable"
	}
	log.Printf("DATABASE_DSN: %s\n", dsn)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	log.Printf("REDIS_ADDR: %s\n", redisAddr)
	redisPwd := os.Getenv("REDIS_PASSWORD")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repo.NewPostgres(ctx, dsn)
	if err != nil {
		log.Fatalf("failed connect pg: %v", err)
	}
	defer db.Close()

	cache := repo.NewRedis(redisAddr, redisPwd)
	defer cache.Close()

	authService := service.NewAuthService(db, cache)

	// start server in goroutine
	go handler.StartServer(authService)

	// graceful shutdown on signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutdown signal received")
}
