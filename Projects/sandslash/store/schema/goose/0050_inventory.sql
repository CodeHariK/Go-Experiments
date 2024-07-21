-- +goose Up

-- Create "inventory" table
CREATE TABLE IF NOT EXISTS "inventory" (
    "id" SERIAL PRIMARY KEY,
    "product_id" INTEGER NOT NULL REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "seller_id" INTEGER NOT NULL REFERENCES "seller" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "quantity" INTEGER NOT NULL
);

CREATE INDEX idx_inventory_seller_id_id ON "inventory" ("seller_id");

CREATE INDEX idx_inventory_product_id ON "inventory" ("product_id");

-- +goose Down
DROP TABLE "inventory";