-- name: CreateUser :one
INSERT INTO
    "users" (
        "username",
        "email",
        "is_admin",
        "date_of_birth",
        "phone_number",
        "last_login",
        "location"
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    "id",
    "username",
    "email",
    "is_admin",
    "created_at",
    "date_of_birth",
    "updated_at",
    "phone_number",
    "last_login",
    "location";

-- name: GetUserByID :one
SELECT
    "id",
    "username",
    "email",
    "is_admin",
    "created_at",
    "date_of_birth",
    "updated_at",
    "phone_number",
    "last_login",
    "location"
FROM "users"
WHERE
    "id" = $1;

-- name: UpdateUser :one
UPDATE "users"
SET
    "username" = $1,
    "email" = $2,
    "is_admin" = $3,
    "date_of_birth" = $4,
    "phone_number" = $5,
    "last_login" = $6,
    "location" = $7,
    "updated_at" = CURRENT_TIMESTAMP
WHERE
    "id" = $8
RETURNING
    "id",
    "username",
    "email",
    "is_admin",
    "created_at",
    "date_of_birth",
    "updated_at",
    "phone_number",
    "last_login",
    "location";

-- name: DeleteUser :one
DELETE FROM "users" WHERE "id" = $1 RETURNING "id";

-- name: ListAllUsers :many
SELECT
    "id",
    "username",
    "email",
    "is_admin",
    "created_at",
    "date_of_birth",
    "updated_at",
    "phone_number",
    "last_login",
    "location"
FROM "users"
ORDER BY "created_at" DESC;

-- name: FindUserByUsername :one
SELECT
    "id",
    "username",
    "email",
    "is_admin",
    "created_at",
    "date_of_birth",
    "updated_at",
    "phone_number",
    "last_login",
    "location"
FROM "users"
WHERE
    "username" = $1;