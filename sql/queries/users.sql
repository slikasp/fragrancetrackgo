-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUser :one
SELECT * 
FROM users 
WHERE name = $1::text;

-- name: GetUsers :many
SELECT name FROM users;

-- name: UpdateUser :one
UPDATE users 
SET name = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetAdmin :one
UPDATE users 
SET is_admin = TRUE, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;