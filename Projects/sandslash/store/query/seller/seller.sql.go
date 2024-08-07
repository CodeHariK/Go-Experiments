// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: seller.sql

package seller

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSeller = `-- name: CreateSeller :one
INSERT INTO seller (name, location) VALUES ($1, $2) RETURNING id
`

type CreateSellerParams struct {
	Name     string      `json:"name"`
	Location pgtype.Int4 `json:"location"`
}

func (q *Queries) CreateSeller(ctx context.Context, arg CreateSellerParams) (int32, error) {
	row := q.db.QueryRow(ctx, createSeller, arg.Name, arg.Location)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteSeller = `-- name: DeleteSeller :exec
DELETE FROM seller WHERE id = $1
`

func (q *Queries) DeleteSeller(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteSeller, id)
	return err
}

const getSellerByID = `-- name: GetSellerByID :one
SELECT id, name, location FROM seller WHERE id = $1
`

func (q *Queries) GetSellerByID(ctx context.Context, id int32) (Seller, error) {
	row := q.db.QueryRow(ctx, getSellerByID, id)
	var i Seller
	err := row.Scan(&i.ID, &i.Name, &i.Location)
	return i, err
}

const listSellers = `-- name: ListSellers :many
SELECT id, name, location FROM seller
`

func (q *Queries) ListSellers(ctx context.Context) ([]Seller, error) {
	rows, err := q.db.Query(ctx, listSellers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Seller{}
	for rows.Next() {
		var i Seller
		if err := rows.Scan(&i.ID, &i.Name, &i.Location); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSeller = `-- name: UpdateSeller :one
UPDATE seller SET name = $1, location = $2 WHERE id = $3 RETURNING id, name, location
`

type UpdateSellerParams struct {
	Name     string      `json:"name"`
	Location pgtype.Int4 `json:"location"`
	ID       int32       `json:"id"`
}

func (q *Queries) UpdateSeller(ctx context.Context, arg UpdateSellerParams) (Seller, error) {
	row := q.db.QueryRow(ctx, updateSeller, arg.Name, arg.Location, arg.ID)
	var i Seller
	err := row.Scan(&i.ID, &i.Name, &i.Location)
	return i, err
}
