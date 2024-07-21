-- +goose Up

-- Create "orders" table
CREATE TABLE IF NOT EXISTS "orders" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "seller_id" INTEGER NOT NULL  REFERENCES "seller" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "status" VARCHAR(12) NOT NULL DEFAULT 'PENDING' CHECK ("status" IN ('PENDING', 'PROCESSING', 'SHIPPED', 'CANCELED', 'REFUNDED')),
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "total_amount" NUMERIC(10, 4) NOT NULL,
    "currency" VARCHAR(12) NOT NULL DEFAULT 'USD' CHECK ("currency" IN ('USD', 'INR', 'EUR'))
);

CREATE INDEX idx_orders_fulfillment_center_id ON "orders" ("seller_id");

CREATE INDEX idx_orders_user_id ON "orders" ("user_id");

-- Create "order_items" table
CREATE TABLE IF NOT EXISTS "order_items" (
    "id" SERIAL PRIMARY KEY,
    "order_id" INTEGER NOT NULL REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "product_id" INTEGER NOT NULL REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "quantity" INTEGER NOT NULL,
    "price" NUMERIC(10, 4) NOT NULL,
    "currency" VARCHAR(12) NOT NULL DEFAULT 'USD' CHECK ("currency" IN ('USD', 'INR', 'EUR'))
);

CREATE INDEX idx_order_items_order_id ON "order_items" ("order_id");

CREATE INDEX idx_order_items_product_id ON "order_items" ("product_id");

-- +goose Down
DROP TABLE "order_items";

DROP TABLE "orders";