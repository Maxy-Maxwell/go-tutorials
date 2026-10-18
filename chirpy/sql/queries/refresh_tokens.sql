-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (token)
DO UPDATE SET
    updated_at = $3,
    expires_at = $5
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshTokenByToken :one
UPDATE refresh_tokens
SET updated_at = $2,
    revoked_at = $3
WHERE token = $1
RETURNING *;
