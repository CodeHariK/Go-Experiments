package services

import (
	"testing"

	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/aggregate"
	"github.com/google/uuid"
)

func init_products(t *testing.T) []aggregate.Product {
	book1, err := aggregate.NewProduct("Book1", "Thinking fast and slow", 1)
	if err != nil {
		t.Fatal(err)
	}

	book2, err := aggregate.NewProduct("Book2", "Human Universe", 2.5)
	if err != nil {
		t.Fatal(err)
	}

	book3, err := aggregate.NewProduct("Book3", "The greates show on earth", 5.2)
	if err != nil {
		t.Fatal(err)
	}

	return []aggregate.Product{
		book1, book2, book3,
	}
}

func TestOrder_NewOrderService(t *testing.T) {
	products := init_products(t)

	os, err := NewOrderService(
		WithMemoryCustomerRepository(),
		WithMemoryProductRepository(products),
	)
	if err != nil {
		t.Fatal(err)
	}

	cust, err := aggregate.NewCustomer("GoExperiments")
	if err != nil {
		t.Error(err)
	}

	err = os.customers.Add(cust)
	if err != nil {
		t.Error(err)
	}

	order := []uuid.UUID{
		products[0].GetID(),
	}

	_, err = os.CreateOrder(cust.GetID(), order)

	if err != nil {
		t.Error(err)
	}
}
