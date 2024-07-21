-- +goose Up

-- Create "product_reviews" table
CREATE TABLE IF NOT EXISTS "product_reviews" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "product_id" INTEGER NOT NULL REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    "rating" INTEGER NOT NULL
);

CREATE INDEX idx_product_reviews_product_id ON "product_reviews" ("product_id");

CREATE INDEX idx_product_reviews_user_id ON "product_reviews" ("user_id");

-- +goose Down
DROP TABLE "product_reviews";