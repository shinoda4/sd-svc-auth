FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
# RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN go env -w GOPROXY=https://goproxy.io
# RUN apk add --no-cache git
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /sd-svc-auth ./cmd/server

FROM alpine:3.18
WORKDIR /app
RUN apk add --no-cache postgresql-client redis curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin
COPY --from=build /sd-svc-auth /app/bin/sd-svc-auth
COPY ./db /app/db
COPY ./scripts /app/scripts
EXPOSE 8080
EXPOSE 50051
ENTRYPOINT ["/app/scripts/run.sh"]
