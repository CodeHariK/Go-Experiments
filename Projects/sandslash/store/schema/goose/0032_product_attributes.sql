-- +goose Up

-- Create the product_attributes table with a variant_id
CREATE TABLE IF NOT EXISTS "product_attributes" (
    "id" SERIAL PRIMARY KEY,
    "product_id" INT NOT NULL REFERENCES "products" ("id") ON DELETE CASCADE,
    "variant_id" INT REFERENCES "product_variants" ("id") ON DELETE CASCADE,
    "attribute_name" VARCHAR(255) NOT NULL,
    "attribute_value" VARCHAR(255) NOT NULL
);

-- Create indexes for the product_attributes table
CREATE INDEX IF NOT EXISTS idx_attribute_product_id ON "product_attributes" ("product_id");

CREATE INDEX IF NOT EXISTS idx_attribute_variant_id ON "product_attributes" ("variant_id");

CREATE INDEX IF NOT EXISTS idx_attribute_name ON "product_attributes" ("attribute_name");

-- +goose Down
DROP TABLE IF EXISTS "product_attributes";