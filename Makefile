include ./.env.example
export


.PHONY: build run docker up init-db docs

check:
	sh scripts/essential.sh

build:
	go build -o bin/sd-svc-auth ./cmd/server

run:
	go run ./cmd/server

deploy-local:
	$(MAKE) build
	sh scripts/run.sh

test:
	go test ./tests/... -v

docker-build:
	docker build --platform=linux/amd64 -t shinoda4/sd-svc-auth:latest .

docker-build-multi:
	docker buildx build \
	--platform linux/amd64,linux/arm64 \
	-t shinoda4/sd-svc-auth:latest \
	--push \
	.

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

kubectl-deploy:
	kubectl delete -f ./manifests ; kubectl apply -f ./manifests