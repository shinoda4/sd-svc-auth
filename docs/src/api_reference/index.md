# API Reference

The API is versioned, currently at `v1`. All endpoints are prefixed with `/api/v1`.

## Base URL

```
http://localhost:8080/api/v1
```

## Content Type

All requests and responses are in JSON format.
`Content-Type: application/json`

## Error Handling

Errors are returned in the following format:

```json
{
  "error": "Description of the error"
}
```

## Authentication

Protected endpoints require a Bearer Token in the Authorization header.

```
Authorization: Bearer <access_token>
```
