-- +goose Up

-- Create "inventory" table
CREATE TABLE IF NOT EXISTS "inventory" (
    "id" SERIAL PRIMARY KEY,
    "product_id" INTEGER NOT NULL,
    "seller_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    CONSTRAINT "inventory_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "inventory_seller_id_fkey" FOREIGN KEY ("seller_id") REFERENCES "seller" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX idx_inventory_seller_id_id ON "inventory" ("seller_id");

CREATE INDEX idx_inventory_product_id ON "inventory" ("product_id");

-- +goose Down
DROP TABLE "inventory";