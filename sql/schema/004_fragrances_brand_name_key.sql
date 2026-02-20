-- +goose Up
ALTER TABLE fragrances
ADD CONSTRAINT fragrances_brand_name_key
UNIQUE (brand, name);

-- +goose Down
ALTER TABLE fragrances
DROP CONSTRAINT fragrances_brand_name_key;