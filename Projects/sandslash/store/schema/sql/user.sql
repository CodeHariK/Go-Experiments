-- name: CreateUser :one
INSERT INTO
    "users" (
        "username",
        "email",
        "is_admin",
        "created_at",
        "date_of_birth",
        "updated_at",
        "phone_number",
        "last_login",
        "address"
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9
    )
RETURNING
    "id";

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
    "address"
FROM "users"
WHERE
    "id" = $1;

-- name: ListUsers :many
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
    "address"
FROM "users";

-- name: UpdateUser :one
UPDATE "users"
SET
    "username" = COALESCE($2, "username"),
    "email" = COALESCE($3, "email"),
    "is_admin" = COALESCE($4, "is_admin"),
    "created_at" = COALESCE($5, "created_at"),
    "date_of_birth" = COALESCE($6, "date_of_birth"),
    "updated_at" = COALESCE($7, "updated_at"),
    "phone_number" = COALESCE($8, "phone_number"),
    "last_login" = COALESCE($9, "last_login"),
    "address" = COALESCE($10, "address")
WHERE
    "id" = $1
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
    "address";

-- name: DeleteUser :exec
DELETE FROM "users" WHERE "id" = $1;