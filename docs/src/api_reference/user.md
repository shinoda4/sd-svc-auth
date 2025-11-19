# User

## Get Current User

Get the profile of the currently authenticated user.

- **URL**: `/authorized/me`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`

### Response

```json
{
  "id": "uuid...",
  "username": "johndoe",
  "email": "user@example.com",
  "email_verified": true
}
```
