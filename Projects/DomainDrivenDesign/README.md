# Domain driven design

[![How To Implement Domain-Driven Design (DDD) in Go](https://img.youtube.com/vi/6zuJXIbOyhs/0.jpg)](https://www.youtube.com/watch?v=6zuJXIbOyhs)

Entity : Unique identifier, Mutable
Value Object : No identifier, Immutable

---

valueobject
    transaction

entity
    item
    person

service
    order
        OrderService struct {customers customer.CustomerRepository, products  product.ProductRepository}
        OrderConfiguration func(OrderService)
        NewOrderService(...OrderConfiguration) OrderService, CreateOrder
    bookstore
        BookStoreConfiguration func(os *BookStore)
        BookStore struct {OrderService, BillingService}
        NewBookStore, WithOrderService, Order(customer uuid.UUID, products []uuid.UUID)

aggregate
    customer
        entity.person, entity.item, valueobject.transactions
        NewCustomer, GetID, SetID, SetName, GetName
    product
        entity.item, price, quantity
        NewProduct, GetID, GetItem

domain
    customer
        repository
            CustomerRepository{Get, Add, Update}
        memory
            memory
                MemoryRepository{customers map[uuid.UUID]aggregate.Customer : New, Get, Add, Update}
        mongo
            mongo
                MongoRepository{Database, Collection : New, Get, Add, Update}
    product
        repository
            ProductRepository{GetAll, GetById, Add, Update, Delete}
        memory
            memory
                MemoryProductRepository{products map[uuid.UUID]aggregate.Product : New, GetAll, GetById, Add, Update, Delete}
