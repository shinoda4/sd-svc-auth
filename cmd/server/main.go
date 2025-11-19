package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shinoda4/sd-svc-auth/internal/config"
	"github.com/shinoda4/sd-svc-auth/internal/repo"
	"github.com/shinoda4/sd-svc-auth/internal/service/auth"
	handler "github.com/shinoda4/sd-svc-auth/internal/transport/http"
	"github.com/shinoda4/sd-svc-auth/pkg/logger"
)

func main() {
	logger.Init()

	// 加载配置（关键）
	cfg := config.MustLoad()

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repo.NewUserRepo(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed connect pg: %v", err)
	}
	defer db.Close()

	cache := repo.NewRedis(cfg.RedisAddr, cfg.RedisPassword)
	defer func(cache *repo.RedisCache) {
		err := cache.Close()
		if err != nil {
			log.Fatalf("failed close redis cache: %v", err)
		}
	}(cache)

	authService := auth.NewAuthService(db, cache)

	// 启动 HTTP 服务
	go handler.StartServer(authService)

	// 优雅关闭
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutdown signal received")
}
