-- +goose Up

-- Create the promotions table
CREATE TABLE IF NOT EXISTS "product_promotions" (
    "id" SERIAL PRIMARY KEY,
    "promotion_name" VARCHAR(255) NOT NULL,
    "discount" NUMERIC(4, 2) NOT NULL,
    "product_variant_id" INT NOT NULL REFERENCES "product_variants" ("id") ON DELETE SET NULL,
    "start_date" DATE NOT NULL,
    "end_date" DATE NOT NULL
);

-- +goose Down
DROP TABLE "product_promotions";