-- name: CreateProduct :one
INSERT INTO
    "products" (
        "product_name",
        "price",
        "category"
    )
VALUES ($1, $2, $3)
RETURNING
    "id";

-- name: GetProductByID :one
SELECT
    "id",
    "product_name",
    "price",
    "category"
FROM "products"
WHERE
    "id" = $1;

-- name: ListProducts :many
SELECT
    "id",
    "product_name",
    "price",
    "category"
FROM "products";

-- name: UpdateProduct :one
UPDATE "products"
SET
    "product_name" = COALESCE($2, "product_name"),
    "price" = COALESCE($3, "price"),
    "category" = COALESCE($4, "category")
WHERE
    "id" = $1
RETURNING
    "id",
    "product_name",
    "price",
    "category";

-- name: DeleteProduct :exec
DELETE FROM "products" WHERE "id" = $1;