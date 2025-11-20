include ./.env.example
export


.PHONY: build run docker up init-db docs

build:
	go build -o bin/sd-svc-auth ./cmd/server

run:
	go run ./cmd/server

deploy:
	$(MAKE) build
	sh scripts/run.sh

test:
	go test ./tests/... -v

docker-build:
	docker build --platform=linux/amd64 -t shinoda4/sd-svc-auth:latest .

docker-up:
	docker-compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

docker-down-v:
	docker compose -f deployments/docker-compose.yml down -v

init-db:
	psql "$(DATABASE_DSN)" -f sql/init.sql -h 127.0.0.1

docs:
	mdbook serve -p 3000 ./docs
