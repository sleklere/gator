-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT feeds.name, url, users.name AS user_name FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;
