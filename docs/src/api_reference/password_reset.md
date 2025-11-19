# Password Reset

## Request Password Reset

Initiate the password reset flow. Sends an email with a reset token.

- **URL**: `/password-reset`
- **Method**: `POST`

### Request Body

```json
{
  "email": "user@example.com",
  "username": "johndoe"
}
```

### Response

```json
{
  "message": "reset email sent" // (Implicit, usually 200 OK)
}
```

## Confirm Password Reset

Reset the password using the token received in the email.

- **URL**: `/password-reset-confirm`
- **Method**: `POST`
- **Query Params**:
  - `token`: The reset token.

### Request Body

```json
{
  "new_password": "newsecurepassword",
  "new_password_confirm": "newsecurepassword"
}
```

### Response

```json
{
  "message": "password updated"
}
```
