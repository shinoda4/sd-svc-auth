# User Profile

The service exposes a single endpoint that returns the authenticated user's profile and token metadata.

## HTTP

- **Method**: `GET /api/v1/me`
- **Headers**: `Authorization: Bearer <access_token>`

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "expires_in": "2025-11-20T12:00:00Z",
  "issued_at": "2025-11-20T11:00:00Z"
}
```

## gRPC

```protobuf
rpc Me(MeRequest) returns (MeResponse);
```

```go
md := metadata.Pairs("authorization", "Bearer "+accessToken)
ctx := metadata.NewOutgoingContext(context.Background(), md)

resp, err := client.Me(ctx, &authpb.MeRequest{})
```

`MeResponse` mirrors the HTTP payload and is derived entirely from the JWT claims that the interceptor injects into `context.Context`.

## Common issues

| Error | Cause | Fix |
|-------|-------|-----|
| `codes.Unauthenticated` / `401` | Missing/expired access token. | Re-login or refresh the session. |
| `codes.PermissionDenied` | Attempting to use a refresh token. | Always send the **access** token to `/me`. |
| 429 / throttling | None built-in. If you see one, a gateway/front-end is rate limiting your requests. | Inspect upstream gateway configuration. |
