-- +goose Up

-- Create "product_reviews" table
CREATE TABLE IF NOT EXISTS "product_reviews" (
    "id" SERIAL PRIMARY KEY,
    "product_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "rating" INTEGER NOT NULL,
    CONSTRAINT "product_reviews_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "product_reviews_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX idx_product_reviews_product_id ON "product_reviews" ("product_id");

CREATE INDEX idx_product_reviews_user_id ON "product_reviews" ("user_id");

-- +goose Down
DROP TABLE "product_reviews";