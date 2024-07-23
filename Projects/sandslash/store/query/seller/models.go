// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package seller

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Author struct {
	ID   int64       `json:"id"`
	Name string      `json:"name"`
	Bio  pgtype.Text `json:"bio"`
}

type GooseDbVersion struct {
	ID        int32            `json:"id"`
	VersionID int64            `json:"version_id"`
	IsApplied bool             `json:"is_applied"`
	Tstamp    pgtype.Timestamp `json:"tstamp"`
}

type Inventory struct {
	ID        int32 `json:"id"`
	ProductID int32 `json:"product_id"`
	SellerID  int32 `json:"seller_id"`
	Quantity  int32 `json:"quantity"`
}

type Location struct {
	ID         int32          `json:"id"`
	Address    string         `json:"address"`
	City       string         `json:"city"`
	State      string         `json:"state"`
	Country    string         `json:"country"`
	PostalCode string         `json:"postal_code"`
	Latitude   pgtype.Numeric `json:"latitude"`
	Longitude  pgtype.Numeric `json:"longitude"`
}

type Order struct {
	ID          int32            `json:"id"`
	UserID      int32            `json:"user_id"`
	Status      string           `json:"status"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	TotalAmount pgtype.Numeric   `json:"total_amount"`
	Currency    string           `json:"currency"`
}

type OrderItem struct {
	ID        int32          `json:"id"`
	OrderID   int32          `json:"order_id"`
	ProductID int32          `json:"product_id"`
	SellerID  int32          `json:"seller_id"`
	Quantity  int32          `json:"quantity"`
	Price     pgtype.Numeric `json:"price"`
	Currency  string         `json:"currency"`
}

type Product struct {
	ID          int32  `json:"id"`
	ProductName string `json:"product_name"`
	Description int32  `json:"description"`
}

type ProductAttribute struct {
	ID             int32       `json:"id"`
	ProductID      int32       `json:"product_id"`
	VariantID      pgtype.Int4 `json:"variant_id"`
	AttributeName  string      `json:"attribute_name"`
	AttributeValue string      `json:"attribute_value"`
}

type ProductComment struct {
	ID        int32            `json:"id"`
	Comment   pgtype.Text      `json:"comment"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type ProductDescription struct {
	ID               int32       `json:"id"`
	ProductID        pgtype.Int4 `json:"product_id"`
	ProductVariantID pgtype.Int4 `json:"product_variant_id"`
	Description      pgtype.Text `json:"description"`
	Images           []string    `json:"images"`
	Videos           []string    `json:"videos"`
}

type ProductPromotion struct {
	ID               int32          `json:"id"`
	PromotionName    string         `json:"promotion_name"`
	Discount         pgtype.Numeric `json:"discount"`
	ProductVariantID int32          `json:"product_variant_id"`
	StartDate        pgtype.Date    `json:"start_date"`
	EndDate          pgtype.Date    `json:"end_date"`
}

type ProductReview struct {
	ID        int32       `json:"id"`
	UserID    int32       `json:"user_id"`
	ProductID int32       `json:"product_id"`
	SellerID  int32       `json:"seller_id"`
	Rating    int32       `json:"rating"`
	Comment   pgtype.Int4 `json:"comment"`
}

type ProductSeller struct {
	ID               int32          `json:"id"`
	ProductVariantID int32          `json:"product_variant_id"`
	SellerID         int32          `json:"seller_id"`
	Price            pgtype.Numeric `json:"price"`
}

type ProductVariant struct {
	ID          int32          `json:"id"`
	ProductID   int32          `json:"product_id"`
	VariantName string         `json:"variant_name"`
	Price       pgtype.Numeric `json:"price"`
	Currency    string         `json:"currency"`
}

type SchemaMigration struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

type Seller struct {
	ID       int32       `json:"id"`
	Name     string      `json:"name"`
	Location pgtype.Int4 `json:"location"`
}

type User struct {
	ID          int32            `json:"id"`
	Username    string           `json:"username"`
	Email       string           `json:"email"`
	PhoneNumber string           `json:"phone_number"`
	IsAdmin     bool             `json:"is_admin"`
	DateOfBirth pgtype.Date      `json:"date_of_birth"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	Location    pgtype.Int4      `json:"location"`
}