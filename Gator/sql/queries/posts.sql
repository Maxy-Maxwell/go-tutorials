-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, published_at, title, url, description, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
ON CONFLICT (url)
DO NOTHING
RETURNING *;

-- name: GetPostsForUser :many
SELECT *
FROM posts
WHERE feed_id = ANY ($1::uuid[])
ORDER BY published_at DESC
LIMIT $2;
