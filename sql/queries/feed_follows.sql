-- Add a CreateFeedFollow query. It will be a deceptively complex SQL query.
-- It should insert a feed follow record, but then return all the fields from the feed follow
-- as well as the names of the linked user and feed.
-- I'll add a tip at the bottom of this lesson if you need it.

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds
ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users
ON users.id = inserted_feed_follow.user_id;

-- It should return all the feed follows for a given user, and include the names of the feeds and user in the result.
-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name, users.name
FROM feed_follows
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
INNER JOIN users
ON users.id = feed_follows.user_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;