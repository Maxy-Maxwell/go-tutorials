-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsWithUserDetails :many
SELECT
    feeds.*,
    users.name AS userName
FROM feeds
LEFT JOIN users
  ON feeds.user_id = users.id;
