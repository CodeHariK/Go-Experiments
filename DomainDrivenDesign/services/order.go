package services

import (
	"context"
	"log"

	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/aggregate"
	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/domain/customer"
	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/domain/customer/memory"
	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/domain/customer/mongo"
	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/domain/product"
	prodmem "github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/domain/product/memory"
	"github.com/google/uuid"
)

type OrderConfiguration func(os *OrderService) error

type OrderService struct {
	customers customer.CustomerRepository
	products  product.ProductRepository
}

func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	os := &OrderService{}
	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

func WithCustomerRepository(cr customer.CustomerRepository) OrderConfiguration {
	return func(os *OrderService) error {
		os.customers = cr
		return nil
	}
}

func WithMongoCustomerRepository(ctx context.Context, connStr string) OrderConfiguration {
	return func(os *OrderService) error {
		cr, err := mongo.New(ctx, connStr)
		if err != nil {
			return err
		}
		os.customers = cr
		return nil
	}
}

func WithMemoryCustomerRepository() OrderConfiguration {
	cr := memory.New()

	return WithCustomerRepository(cr)
}

func WithMemoryProductRepository(products []aggregate.Product) OrderConfiguration {
	return func(os *OrderService) error {
		pr := prodmem.New()

		for _, p := range products {
			if err := pr.Add(p); err != nil {
				return err
			}
		}

		os.products = pr
		return nil
	}
}

func (o *OrderService) CreateOrder(customerID uuid.UUID, productsIDs []uuid.UUID) (float64, error) {
	c, err := o.customers.Get(customerID)
	if err != nil {
		return 0, err
	}

	var products []aggregate.Product
	var total float64

	for _, id := range productsIDs {
		p, err := o.products.GetById(id)
		if err != nil {
			return 0, err
		}

		products = append(products, p)
		total += p.GetPrice()
	}

	log.Printf("Customer: %s has ordered %d products", c.GetID(), len(products))

	return total, nil
}
