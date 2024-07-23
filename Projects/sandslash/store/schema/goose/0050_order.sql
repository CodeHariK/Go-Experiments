-- +goose Up

-- Create "orders" table
CREATE TABLE IF NOT EXISTS "orders" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "status" VARCHAR(12) NOT NULL DEFAULT 'PENDING' CHECK (
        "status" IN (
            'PENDING',
            'PROCESSING',
            'SHIPPED',
            'CANCELED',
            'REFUNDED'
        )
    ),
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "total_amount" NUMERIC(10, 4) NOT NULL,
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

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON "orders" ("user_id");

-- +goose Down

DROP TABLE IF EXISTS "orders";