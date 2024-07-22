-- +goose Up

-- Create the products table
CREATE TABLE IF NOT EXISTS "products" (
    "id" SERIAL PRIMARY KEY,
    "product_name" VARCHAR(255) NOT NULL UNIQUE,
    "description" VARCHAR(2048)
);

-- Create indexes for the products table
CREATE INDEX idx_product_name ON "products" ("product_name");

-- +goose Down
DROP TABLE "products";