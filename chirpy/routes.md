# Chirpy API routes

Base URL: `http://localhost:8080`

Timestamps in JSON bodies are RFC 3339 UTC values, for example `2026-10-19T04:33:00Z`.

Routes that use JWT authentication expect:

```
Authorization: Bearer <access_token>
```

The access token is the `token` field returned by login or refresh. It expires one hour after it is issued.

---

## `GET /app/`

Serves static files from the `./app` directory. The `/app/` prefix is stripped before the file is looked up, so `GET /app/index.html` serves `./app/index.html`. Each request increments the file-server hit counter used by `GET /admin/metrics`.

**Request body:** none

**Success response:** `200 OK` with the file contents and the content type inferred from the file.

---

## Users

### `POST /api/users`

Creates a user.

**Request body:**

```json
{
  "email": "lane@example.com",
  "password": "supersecret"
}
```

**Success response:** `201 Created`

`Content-Type: text/html; charset=utf-8`

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "created_at": "2026-10-19T04:33:00Z",
  "updated_at": "2026-10-19T04:33:00Z",
  "email": "lane@example.com",
  "is_chirpy_red": false
}
```

### `PUT /api/users`

Updates the authenticated user's email and password. Requires a JWT access token.

**Request body:**

```json
{
  "email": "lane.new@example.com",
  "password": "newsecret"
}
```

**Success response:** `200 OK`

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "created_at": "2026-10-19T04:33:00Z",
  "updated_at": "2026-10-19T05:10:00Z",
  "email": "lane.new@example.com",
  "is_chirpy_red": false
}
```

### `POST /api/polka/webhooks`

Polka membership webhook. Requires the Polka API key:

```
Authorization: ApiKey <polka_apikey>
```

**Request body** for a Chirpy Red upgrade:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7"
  }
}
```

**Success response:** `204 No Content`

```json
""
```

Any other `event` value is also accepted and returns the same `204 No Content` body. The user's Chirpy Red status is updated only when `event` is `user.upgraded`.

---

## Login and tokens

### `POST /api/login`

Authenticates a user and issues an access token and a refresh token. The refresh token expires 60 days after it is issued.

**Request body:**

```json
{
  "email": "lane@example.com",
  "password": "supersecret"
}
```

**Success response:** `200 OK`

`Content-Type: text/html; charset=utf-8`

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "created_at": "2026-10-19T04:33:00Z",
  "updated_at": "2026-10-19T04:33:00Z",
  "email": "lane@example.com",
  "token": "<jwt_access_token>",
  "refresh_token": "<refresh_token>",
  "is_chirpy_red": false
}
```

### `POST /api/refresh`

Issues a new access token from a valid refresh token. Send the refresh token as a bearer token:

```
Authorization: Bearer <refresh_token>
```

**Request body:** none

**Success response:** `200 OK`

`Content-Type: text/html; charset=utf-8`

```json
{
  "token": "<jwt_access_token>"
}
```

### `POST /api/revoke`

Revokes a refresh token. Send the refresh token as a bearer token:

```
Authorization: Bearer <refresh_token>
```

**Request body:** none

**Success response:** `204 No Content`

`Content-Type: text/html; charset=utf-8`

The body is empty.

---

## Chirps

### `POST /api/chirps`

Creates a chirp for the authenticated user. Requires a JWT access token. `body` must be 140 characters or fewer. The words `kerfuffle`, `sharbert`, and `fornax` are replaced with `****` (matched case-insensitively, per word).

**Request body:**

```json
{
  "body": "This is a kerfuffle"
}
```

**Success response:** `201 Created`

```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "created_at": "2026-10-19T04:40:00Z",
  "updated_at": "2026-10-19T04:40:00Z",
  "body": "This is a ****",
  "user_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7"
}
```

### `GET /api/chirps`

Returns chirps. An optional `author_id` query parameter limits the list to one user:

```
GET /api/chirps?author_id=7c9e6679-7425-40de-944b-e07fc1f90ae7
```

**Request body:** none

**Success response:** `200 OK`

```json
[
  {
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "created_at": "2026-10-19T04:40:00Z",
    "updated_at": "2026-10-19T04:40:00Z",
    "body": "This is a ****",
    "user_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7"
  }
]
```

When the query returns no chirps, the body is JSON `null`.

### `GET /api/chirps/{id}`

Returns one chirp. `{id}` is the chirp UUID.

**Request body:** none

**Success response:** `200 OK`

```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "created_at": "2026-10-19T04:40:00Z",
  "updated_at": "2026-10-19T04:40:00Z",
  "body": "This is a ****",
  "user_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7"
}
```

### `DELETE /api/chirps/{id}`

Deletes a chirp owned by the authenticated user. Requires a JWT access token. `{id}` is the chirp UUID.

**Request body:** none

**Success response:** `204 No Content`

`Content-Type: text/html; charset=utf-8`

The body is empty.

---

## Health and admin

### `GET /api/healthz`

Readiness check.

**Request body:** none

**Success response:** `200 OK`

`Content-Type: text/plain; charset=utf-8`

```
OK
```

### `GET /admin/metrics`

Returns an HTML page with the current file-server hit count.

**Request body:** none

**Success response:** `200 OK`

`Content-Type: text/html; charset=utf-8`

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 12 times!</p>
  </body>
</html>
```

### `POST /admin/reset`

Resets the file-server hit counter to `0` and deletes every user. Succeeds when `PLATFORM` is `dev`.

**Request body:** none

**Success response:** `200 OK`

`Content-Type: text/plain; charset=utf-8`

```
Hits: 0
```
