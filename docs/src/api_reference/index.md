# API Reference

The authentication service provides both **native gRPC** and **HTTP/JSON** APIs. Choose the protocol that best fits your client requirements.

## Protocol Options

### gRPC API

Native gRPC service using Protocol Buffers for efficient, type-safe communication.

- **Best for**: Microservices, Go/Java/Python clients, high-performance applications
- **Protocol**: HTTP/2 with Protocol Buffers
- **Port**: `50051` (default)
- **Documentation**: [gRPC API Reference](./grpc.md)

### HTTP Gateway API

RESTful-style HTTP/JSON API automatically generated from the gRPC service using grpc-gateway.

- **Best for**: Web browsers, JavaScript clients, REST-based integrations
- **Protocol**: HTTP/1.1 or HTTP/2 with JSON
- **Port**: `8080` (default)
- **Documentation**: [HTTP Gateway API Reference](./http_gateway.md)

## Base URLs

**gRPC**:
```
localhost:50051
```

**HTTP Gateway**:
```
http://localhost:8080
```

## Content Type

**gRPC**: Protocol Buffers (binary)

**HTTP**: JSON
```
Content-Type: application/json
```

## Authentication

Both protocols support Bearer token authentication.

### gRPC Metadata

```
authorization: Bearer <access_token>
```

### HTTP Header

```
Authorization: Bearer <access_token>
```

## Error Handling

### gRPC Status Codes

The service uses standard gRPC status codes:

- `OK` (0): Success
- `INVALID_ARGUMENT` (3): Invalid input
- `UNAUTHENTICATED` (16): Authentication failed
- `ALREADY_EXISTS` (6): Resource already exists
- `NOT_FOUND` (5): Resource not found
- `INTERNAL` (13): Server error

### HTTP Status Codes

The HTTP gateway automatically maps gRPC codes to HTTP status codes:

- `200 OK`: Success
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication failed
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

## API Sections

- **[gRPC API](./grpc.md)**: Native gRPC service methods and examples
- **[HTTP Gateway API](./http_gateway.md)**: HTTP/JSON endpoints and examples
- **[Authentication](./auth.md)**: Deprecated (see gRPC/HTTP Gateway docs)
- **[User](./user.md)**: Deprecated (see gRPC/HTTP Gateway docs)
- **[Password Reset](./password_reset.md)**: Deprecated (see gRPC/HTTP Gateway docs)

