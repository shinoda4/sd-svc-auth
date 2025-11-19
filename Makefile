include ./.env.example
export


.PHONY: build run docker up init-db

build:
	go build -o bin/sd-svc-auth ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./tests/... -v

docker:
	docker build -t sd-svc-auth:local .

up:
	docker-compose -f deployments/docker-compose.yml up -d

init-db:
	# run init SQL (requires psql installed)
	psql "$(DATABASE_DSN)" -f sql/init.sql -h 127.0.0.1
