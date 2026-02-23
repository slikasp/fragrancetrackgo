-- name: AddRating :one
INSERT INTO ratings (user_id, brand, name, rating, comment)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetRating :one
SELECT * 
FROM ratings 
WHERE user_id = $1 AND brand = $2 AND name = $3;

-- name: GetRatings :many
SELECT * FROM ratings
WHERE user_id = $1;

-- name: UpdateRating :one
UPDATE ratings 
SET rating = $4, comment = $5
WHERE user_id = $1 AND brand = $2 AND name = $3
RETURNING *;

-- name: RemoveRating :one
DELETE FROM ratings 
WHERE user_id = $1 AND brand = $2 AND name = $3
RETURNING *;

-- name: ResetRatings :exec
DELETE FROM ratings;