-- +goose Up

-- Create "products" table
CREATE TABLE IF NOT EXISTS "products" (
    "id" SERIAL PRIMARY KEY,
    "product_name" VARCHAR(255) NOT NULL UNIQUE,
    "category" VARCHAR(255) NOT NULL,
    "price" NUMERIC(10, 4) NOT NULL,
    "currency" VARCHAR(12) NOT NULL DEFAULT 'USD' CHECK ("currency" IN ('USD', 'INR', 'EUR'))
);

CREATE INDEX idx_products_category_id ON "products" ("category");

-- +goose Down
DROP TABLE "products";