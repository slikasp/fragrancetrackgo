-- name: AddFragrance :one
INSERT INTO fragrances (brand, name)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetFragrance :one
SELECT * 
FROM fragrances 
WHERE brand = $1 AND name = $2;

-- name: GetFragrances :many
SELECT * FROM fragrances;

-- name: UpdateFragrance :one
UPDATE fragrances 
SET brand = $2, name = $3
WHERE id = $1
RETURNING *;

-- name: RemoveFragrance :one
DELETE FROM fragrances WHERE id = $1
RETURNING *;

-- name: ResetFragrances :exec
DELETE FROM fragrances;