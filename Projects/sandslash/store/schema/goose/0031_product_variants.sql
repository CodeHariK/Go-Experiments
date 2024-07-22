-- +goose Up

-- Create the product_variants table
CREATE TABLE IF NOT EXISTS "product_variants" (
    "id" SERIAL PRIMARY KEY,
    "product_id" INT NOT NULL REFERENCES "products" ("id") ON DELETE CASCADE,
    "variant_name" VARCHAR(255) NOT NULL,
    "images" VARCHAR(1024) [],
    "videos" VARCHAR(1024) [],
    "description" VARCHAR(1024),
    "price" NUMERIC(10, 4) NOT NULL,
    "currency" VARCHAR(12) NOT NULL DEFAULT 'USD' CHECK (
        "currency" IN (
            'USD',
            'INR',
            'BTC',
            'ETH',
            'SOL'
        )
    )
);

-- Create indexes for the product_variants table
CREATE INDEX idx_variant_product_id ON "product_variants" ("product_id");

CREATE INDEX idx_variant_name ON "product_variants" ("variant_name");

-- +goose Down
DROP TABLE "product_variants";