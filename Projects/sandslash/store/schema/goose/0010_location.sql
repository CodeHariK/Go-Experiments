-- +goose Up

-- Create "locations" table
CREATE TABLE IF NOT EXISTS "locations" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "city" VARCHAR(100) NOT NULL,
    "state" VARCHAR(100) NOT NULL,
    "country" VARCHAR(100) NOT NULL,
    "postal_code" VARCHAR(20) NOT NULL,
    "latitude" NUMERIC(9, 6) NOT NULL,
    "longitude" NUMERIC(9, 6) NOT NULL
);

-- Create an index on the latitude and longitude for efficient spatial queries
CREATE INDEX IF NOT EXISTS idx_locations_latitude_longitude ON "locations" ("latitude", "longitude");

-- +goose Down

DROP TABLE IF EXISTS "locations";