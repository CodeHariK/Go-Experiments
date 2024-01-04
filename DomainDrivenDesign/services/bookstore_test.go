package services

import (
	"context"
	"testing"

	"github.com/CodeHariK/Go-Experiments/DomainDrivenDesign/aggregate"
	"github.com/google/uuid"
)

func Test_BookStore(t *testing.T) {
	products := init_products(t)

	os, err := NewOrderService(
		//__________________
		// WithMemoryCustomerRepository(),
		WithMongoCustomerRepository(context.Background(), "mongodb+srv://admin:GoExperiments@cluster0.tendsox.mongodb.net/?retryWrites=true&w=majority"),
		//__________________

		WithMemoryProductRepository(products),
	)
	if err != nil {
		t.Fatal(err)
	}

	bookstore, err := NewBookStore(WithOrderService(os))
	if err != nil {
		t.Fatal(err)
	}

	cust, err := aggregate.NewCustomer("GoExperiments")
	if err != nil {
		t.Fatal(err)
	}

	if err = os.customers.Add(cust); err != nil {
		t.Fatal(err)
	}

	order := []uuid.UUID{
		products[0].GetID(),
	}

	err = bookstore.Order(cust.GetID(), order)
	if err != nil {
		t.Fatal(err)
	}
}
