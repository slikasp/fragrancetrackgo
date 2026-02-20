-- +goose Up
INSERT INTO users(id, created_at, updated_at, name, hashed_password, is_admin)
VALUES (
    '9814a5c0-78e1-4aea-be47-8fdff0385bb5',
    NOW(),
    NOW(),
    'slip',
    '$argon2id$v=19$m=65536,t=3,p=2$8geVDR9aTIayNiHxPGy5yA$sqprOBBo9gLJmJVDKxw8glPTQNHHe8mcufzcTIxt1d4',
    't'
);

-- +goose Down
DELETE FROM users WHERE id = '9814a5c0-78e1-4aea-be47-8fdff0385bb5';