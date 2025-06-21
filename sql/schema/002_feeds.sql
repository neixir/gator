-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    -- https://stackoverflow.com/questions/14141266/postgresql-foreign-key-on-delete-cascade/14141354#14141354
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL
    -- o tambe podriem dir
    -- user_id UUID,
    -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;