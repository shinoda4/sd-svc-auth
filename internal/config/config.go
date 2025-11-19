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
		"SERVER_PORT",
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
		ServerPort:    os.Getenv("SERVER_PORT"),
		GrpcPort:      os.Getenv("GRPC_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		EmailAddress:  os.Getenv("EMAIL_ADDRESS"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}
