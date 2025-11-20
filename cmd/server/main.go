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
	"github.com/shinoda4/sd-svc-auth/internal/transport/grpc"
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

	go grpc.RunGRPCServer(authService) // gRPC server
	go grpc.RunGateway()
	//go handler.StartServer(authService) // Http server

	// 优雅关闭
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutdown signal received")
}
