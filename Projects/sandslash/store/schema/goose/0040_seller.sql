-- +goose Up

-- Create "seller" table
CREATE TABLE IF NOT EXISTS "seller" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "location" VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE "seller";