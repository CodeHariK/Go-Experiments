-- +goose Up

-- Create "orders" table
CREATE TABLE IF NOT EXISTS "orders" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "seller_id" INTEGER NOT NULL,
    "total_amount" NUMERIC(10, 2) NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "orders_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "orders_seller_id_fkey" FOREIGN KEY ("seller_id") REFERENCES "seller" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX idx_orders_fulfillment_center_id ON "orders" ("seller_id");

CREATE INDEX idx_orders_user_id ON "orders" ("user_id");

-- Create "order_items" table
CREATE TABLE IF NOT EXISTS "order_items" (
    "id" SERIAL PRIMARY KEY,
    "order_id" INTEGER NOT NULL,
    "product_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "price" NUMERIC(10, 2) NOT NULL,
    CONSTRAINT "order_items_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "order_items_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX idx_order_items_order_id ON "order_items" ("order_id");

CREATE INDEX idx_order_items_product_id ON "order_items" ("product_id");

-- +goose Down
DROP TABLE "order_items";

DROP TABLE "orders";