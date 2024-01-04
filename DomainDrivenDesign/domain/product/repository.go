package product

import (
	"errors"

	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/aggregate"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound      = errors.New("no such product")
	ErrProductAlreadyExists = errors.New("there is already such a product")
)

type ProductRepository interface {
	GetAll() ([]aggregate.Product, error)
	GetById(id uuid.UUID) (aggregate.Product, error)
	Add(product aggregate.Product) error
	Update(product aggregate.Product) error
	Delete(id uuid.UUID) error
}
