module github.com/shinoda4/sd-svc-auth

go 1.25.4

// replace github.com/shinoda4/sd-grpc-proto => ../sd-grpc-proto

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.16.0
	github.com/shinoda4/sd-grpc-proto v0.0.7
	golang.org/x/crypto v0.44.0
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251111163417-95abcf5c77ba // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251111163417-95abcf5c77ba // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)
