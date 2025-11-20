#!/bin/sh
set -e

# detect environment
if ping -c 1 postgres >/dev/null 2>&1; then
  PG_HOST=postgres
else
  PG_HOST=127.0.0.1
fi

if ping -c 1 redis >/dev/null 2>&1; then
  REDIS_HOST=redis
else
  REDIS_HOST=127.0.0.1
fi

echo "Waiting for Postgres..."
until pg_isready -h $PG_HOST -p 5432 -U sd_svc_auth_user; do
  echo "Postgres not ready, sleeping..."
  sleep 1
done

echo "Waiting for Redis..."
until redis-cli -h $REDIS_HOST -p 6379 ping | grep -q PONG; do
  echo "Redis not ready, sleeping..."
  sleep 1
done
