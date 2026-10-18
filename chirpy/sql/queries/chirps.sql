-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: DeleteChirps :exec
DELETE FROM chirps;

-- name: DeleteChirpById :exec
DELETE FROM chirps
WHERE
    id = $1
    AND user_id = $2;


-- name: GetChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT *
FROM chirps
WHERE id = $1;
