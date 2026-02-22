-- +goose Up
ALTER TABLE ratings
ADD CONSTRAINT ratings_key
UNIQUE (user_id, brand, name);

-- +goose Down
ALTER TABLE ratings
DROP CONSTRAINT ratings_key;