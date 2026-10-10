-- name: GetFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted_feed_follow.*, users.name, feeds.name
FROM inserted_feed_follow
LEFT JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
LEFT JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT *, feeds.name AS feeds_name, users.name AS users_name
FROM feed_follows
LEFT JOIN feeds ON feed_follows.feed_id = feeds.id
LEFT JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollowsForUserByFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;
