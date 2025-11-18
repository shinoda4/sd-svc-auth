# API Reference

Prefix: `/api/v1`

[[toc]]

## Public

Prefix: `/`

### Register

Create a new user.

Parameters:

`sendEmail`: `true` or `false`

```text
POST /api/v1/register
{
    "email": "user@example.com",
    "username": "hijack",
    "password": "123456"
}
```

Return:

```text
201 Created
{
    "message": "registered"
}
```

### Login

Login user and returns JWT token.

```text
POST /api/v1/login
{
    "email": "user@example.com",
    "password": "******"
}
```

Return:

```text
200 OK
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJlMWJkZDE0YS1lMDgzLTRkY2UtYjc0OC04MDFlZWFiNTQzNzMiLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NjMxMzM5MjAsImlhdCI6MTc2MzEzMDMyMH0.X-GCvWWWXa3imkWloe_QJoeYTJhQUFVLHMPR4W6Uy_o",
    "expires_in": 3600,
    "refresh_expires_in": 259200,
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJlMWJkZDE0YS1lMDgzLTRkY2UtYjc0OC04MDFlZWFiNTQzNzMiLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NjMzODk1MjAsImlhdCI6MTc2MzEzMDMyMH0.HQNwXVwpEzERQSTNdw30Ae41sgZ1tTNm7L9tsf8I6vY"
}
```

### Refresh

Refresh 

```text
POST /api/v1/refresh
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJlMWJkZDE0YS1lMDgzLTRkY2UtYjc0OC04MDFlZWFiNTQzNzMiLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NjMzODk1MjAsImlhdCI6MTc2MzEzMDMyMH0.HQNwXVwpEzERQSTNdw30Ae41sgZ1tTNm7L9tsf8I6vY"
}
```

Return:

```text
200 OK
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJlMWJkZDE0YS1lMDgzLTRkY2UtYjc0OC04MDFlZWFiNTQzNzMiLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3NjMxMzQwMDcsImlhdCI6MTc2MzEzMDQwN30.NlF9yGojW9gdQ3VIOpMv1rVpJ0vy9ODtlhiyko1Da9E",
    "expires_in": 3600
}
```

### Verify

Verify token is valid.

```text
POST /api/v1/verify
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwidWlkIjoiZTFiZGQxNGEtZTA4My00ZGNlLWI3NDgtODAxZWVhYjU0MzczIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzYzMTM4MDc1LCJpYXQiOjE3NjMxMzgwNzR9.pzVDuKoyUu5EAv-E_vi-fhitSEruqNUZhmNduNzLsms"
}
```

Return:

```text
200 OK
{
    "token": {
        "token_type": "access",
        "uid": "e1bdd14a-e083-4dce-b748-801eeab54373",
        "email": "user@example.com",
        "exp": 1763141504,
        "iat": 1763137904
    }
}
```

### Password Reset

Required for password reset.

```text
POST /api/v1/password-reset

{
	"email": "lindesong666@163.com",
    "username": "desonglll"
}
```

You should config the `RESET_PASSWORD_URL` environment variable.

This operation will create a new token and expire time, then insert it to `reset_token` field and `reset_token_expire` in `users` table. The default value of these two field is `null`.

When the reset process is done successfully, these two fields shoud back to `null`.

Finally, the system will send an email that contains a link which should redirect to `password reset confirm` api to the required email address. You can change the confirmation host by set `RESET_PASSWORD_URL`.

The confirmation link should looks like `RESET_PASSWORD_URL?token=TOKEN`.

Return:

```text
200 OK
```

### Password Reset Confirm

Confirm the password reset action.

Before access this api, you should already access the `password-reset` api successfully.

In the past section, we described that the system is send an email that contains confirmation link to required email address.

```text
POST api/v1/password-reset-confirm?token=[TOKEN]

{
    "new_password": "123456",
    "new_password_confirm": "123456"
}
```

Return:

```text
200 OK
```

## Authorized

Prefix: `/authorized`

### Me

```text
GET /api/v1/authorized/me
Authorization: Bearer xxx
```

Return:

```text
{
    "claims": {
        "uid": "e1bdd14a-e083-4dce-b748-801eeab54373",
        "email": "user@example.com",
        "exp": 1763386825,
        "iat": 1763127625
    }
}
```

### Logout

Logout user.

```text
POST /api/v1/logout
```

Return:

```text
200 OK
{
    "message": "logout successful"
}
```

To Be Done...
