#!/bin/sh
set -e

# 在 Kubernetes 环境中，PG_HOST 和 REDIS_HOST 应该始终是 Service 名称
PG_HOST=postgres
REDIS_HOST=redis
REDIS_PASSWORD="${REDIS_PASSWORD}"

echo "Waiting for Postgres at $PG_HOST..."
# 使用 pg_isready 检查连接，这是最可靠的检测方式
until pg_isready -h "$PG_HOST" -p 5432 -U sd_svc_auth_user; do
  echo "Postgres not ready, sleeping..."
  sleep 1
done

echo "Waiting for Redis at $REDIS_HOST..."
# 使用 redis-cli ping 检查连接
until redis-cli -h "$REDIS_HOST" -p 6379 -a "$REDIS_PASSWORD" ping | grep -q PONG; do
  echo "Redis not ready, sleeping..."
  sleep 1
done

echo "Database checks passed. Starting application..."
