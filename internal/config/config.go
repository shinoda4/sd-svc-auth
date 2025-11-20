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

package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseDSN   string
	RedisAddr     string
	RedisPassword string
	ServerHost    string
	ServerPort    string
	GrpcPort      string
	JWTSecret     string
	EmailAddress  string
	EmailPassword string
}

func MustLoad() *Config {
	required := []string{
		"DATABASE_DSN",
		"REDIS_ADDR",
		"SERVER_HOST",
		"GRPC_PORT",
		"JWT_SECRET",
		"EMAIL_ADDRESS",
		"EMAIL_PASSWORD",
	}

	// 统一检查缺失变量
	var missing []string
	for _, env := range required {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("missing required environment variables: %v", missing)
	}

	return &Config{
		DatabaseDSN:   os.Getenv("DATABASE_DSN"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"), // 可选
		ServerHost:    os.Getenv("SERVER_HOST"),
		GrpcPort:      os.Getenv("GRPC_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		EmailAddress:  os.Getenv("EMAIL_ADDRESS"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}
