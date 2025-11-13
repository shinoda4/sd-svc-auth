FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN apk add --no-cache git
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /sd-svc-auth ./cmd/server

FROM scratch
COPY --from=build /sd-svc-auth /sd-svc-auth
EXPOSE 8080
ENTRYPOINT ["/sd-svc-auth"]
