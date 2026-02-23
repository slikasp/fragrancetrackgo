-- name: SearchFragrances :many
SELECT id, brand, name, url
FROM public.fragrances
WHERE (
  SELECT bool_and( (brand || ' ' || name) ILIKE '%' || t || '%' )
  FROM unnest(regexp_split_to_array(trim($1), '\s+')) AS t
)
ORDER BY brand, name
LIMIT $2
OFFSET $3;