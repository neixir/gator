-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE NOT NULL,
    FOREIGN KEY (user_id, feed_id),
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;