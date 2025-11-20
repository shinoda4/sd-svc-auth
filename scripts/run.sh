#!/bin/sh

sh ./scripts/essential.sh

sh ./scripts/wait-for-deps.sh

sh ./scripts/migrate.sh

echo "All dependencies are ready, starting auth-service..."

exec ./bin/sd-svc-auth