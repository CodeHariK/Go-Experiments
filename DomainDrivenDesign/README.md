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
    bookstore
    order

aggregate
    customer
        person, products, transactions
        NewCustomer, GetID, SetID, SetName, GetName
    product

domain
    customer
        repository
        memory
            memory
        mongo
            mongo
    product
        repository
        memory
            memory
