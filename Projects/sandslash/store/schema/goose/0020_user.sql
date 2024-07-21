-- +goose Up

-- Create "users" table
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(50) NOT NULL UNIQUE,
    "email" VARCHAR(50) NOT NULL UNIQUE,
    "is_admin" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "date_of_birth" DATE,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "phone_number" VARCHAR(15) NOT NULL,
    "last_login" TIMESTAMP,
    "address" VARCHAR(255)
);

CREATE INDEX "idx_users_email" ON "users" ("email");

-- +goose Down

DROP TABLE "users";