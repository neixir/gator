-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url=$1;

-- CH5 L1
-- It should simply set the last_fetched_at and updated_at columns to the current time for a given feed (probably by ID is simplest).
-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1;

-- CH5 L1
-- It should return the next feed we should fetch posts from.
-- We want to scrape all the feeds in a continuous loop.
-- A simple approach is to keep track of when a feed was last fetched,
-- and always fetch the oldest one first (or any that haven't ever been fetched).
-- SQL has a NULLS FIRST clause that can help with this.
-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;