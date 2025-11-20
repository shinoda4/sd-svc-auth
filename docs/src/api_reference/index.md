# API Reference

`sd-svc-auth` exposes the same surface area over two transports:

- **gRPC** – strongly typed Protocol Buffers served on `$GRPC_PORT` (default `50051`).
- **HTTP/JSON** – automatically generated REST endpoints served on `$HTTP_PORT` (default `8080`) under `/api/v1/*`.

The protobuf definitions live in [`github.com/shinoda4/sd-grpc-proto/proto/auth/v1`](https://github.com/shinoda4/sd-grpc-proto). The grpc-gateway uses the `google.api.http` annotations from that repository to produce HTTP routes.

## Base endpoints

| Transport | Address | Notes |
|-----------|---------|-------|
| gRPC | `localhost:50051` | Use TLS for production deployments. |
| HTTP | `http://localhost:8080/api/v1` | Prefix is already included in every example below. |

## Authentication

Protected methods require a Bearer token. Registration, login, verification, forgot/reset password, and health check remain public.

- **gRPC metadata**
  ```
  authorization: Bearer <token>
  ```
- **HTTP header**
  ```
  Authorization: Bearer <token>
  ```

When a method expects a refresh token, send it through the same header (the interceptor determines token type automatically).

## Error model

| gRPC Status | HTTP Status | Typical cause |
|-------------|-------------|---------------|
| `INVALID_ARGUMENT` | `400 Bad Request` | Missing fields, password mismatch, malformed token. |
| `UNAUTHENTICATED` | `401 Unauthorized` | Missing/invalid Bearer token. |
| `NOT_FOUND` | `404 Not Found` | Unknown tokens or user records. |
| `ALREADY_EXISTS` | `409 Conflict` | Duplicate email or username. |
| `INTERNAL` | `500 Internal Server Error` | Unexpected server/database issue. |

## Reference sections

- [gRPC API](./grpc.md) – Proto method signatures and language-specific snippets.
- [HTTP Gateway API](./http_gateway.md) – JSON endpoints, curl samples, and request/response payloads.
- [Authentication & Tokens](./auth.md) – Conceptual guide to registration, login, refresh, and logout flows.
- [User Profile](./user.md) – How to retrieve the current user.
- [Password Reset Flow](./password_reset.md) – Requesting and confirming password resets through both transports.
