-- +goose Up

-- Create "products" table
CREATE TABLE IF NOT EXISTS "products" (
    "id" SERIAL PRIMARY KEY,
    "product_name" VARCHAR(255) NOT NULL UNIQUE,
    "price" NUMERIC(10, 2) NOT NULL,
    "category" VARCHAR(255) NOT NULL
);

CREATE INDEX idx_products_category_id ON "products" ("category");

-- +goose Down
DROP TABLE "products";