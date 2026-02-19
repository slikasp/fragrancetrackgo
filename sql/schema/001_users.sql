-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NUll,
    hashed_password TEXT NOT NULL DEFAULT 'unset',
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE users;