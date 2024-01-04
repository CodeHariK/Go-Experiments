package services

import (
	"log"

	"github.com/google/uuid"
)

type BookStoreConfiguration func(os *BookStore) error

type BookStore struct {
	OrderService *OrderService

	BillingService interface{}
}

func NewBookStore(cfgs ...BookStoreConfiguration) (*BookStore, error) {
	t := &BookStore{}

	for _, cfg := range cfgs {
		if err := cfg(t); err != nil {
			return nil, err
		}
	}
	return t, nil
}

func WithOrderService(os *OrderService) BookStoreConfiguration {
	return func(t *BookStore) error {
		t.OrderService = os
		return nil
	}
}

func (t *BookStore) Order(customer uuid.UUID, products []uuid.UUID) error {
	price, err := t.OrderService.CreateOrder(customer, products)
	if err != nil {
		return err
	}

	log.Printf("\nBill the customer: %0.0f\n", price)

	return nil
}
