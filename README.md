
Cart API
========

## Overview

This repository contains a Go programming exercise for interview candidates.  
You'll be developing an API for an online shopping cart in the Go programming language.


## Requirements

This is a REST API for basic CRUD operations for an online shopping cart. Data
should be persisted in a storage layer which can use Postgres.

You should use the default `net/http` package for REST implementation; `sqlx` or `sqlc` for interacting with Postgres;
all the queries should be written manually (no ORM); your repo should be private.

#### Additional requirements

- Cover your code with unit tests (you could use `testify`).
- Create a `Dockerfile` and `Docker compose` for your application.
- For parsing your environment files use `viper`.
- To apply migrations for Postgres use `goose` in your application.

### Domain Types

The Cart API consists of two simple types: `Cart` and `CartItem`. The `Cart`  
holds zero or more `CartItem` objects.

`CartItem` objects should be created in DB exactly (not from application).

- The maximum number of distinct products in one cart is **5**.

### Create Cart

A new cart should be created and an ID generated. The new empty cart should be returned.

```sh
POST http://localhost:3000/carts -d '{}'
```

```json
{
  "id": 1,
  "items": []
}
```

### Add to Cart

Cart can contain only 5 products

A new item should be added to an existing cart. 
The new item should be returned.

Should fail if:
  - The cart does not exist.
  - The product name is blank.
  - The price is non-positive.
  - The cart already contains 5 products.

```sh
POST http://localhost:3000/carts/1/items -d '{
  "product": "Shoes",
  "price": 2500.50
}'
```

```json
{
  "id": 1,
  "cart_id": 1,
  "product": "Shoes",
  "price": 2500.50
}
```

### Remove from Cart

An existing item should be removed from a cart. Should fail if the cart does not
exist or if the item does not exist.

```sh
DELETE http://localhost:3000/carts/1/items/1
```

```json
{}
```


### View Cart

An existing cart should be able to be viewed with its items. Should fail if the
cart does not exist.

```sh
GET http://localhost:3000/carts/1
```

```json
{
  "id": 1,
  "items": [
    {
      "id": 1,
      "cart_id": 1,
      "product": "Shoes",
      "price": 2500.50
    },
    {
      "id": 2,
      "cart_id": 1,
      "product": "Socks",
      "price": 1200.00
    }
  ]
}
```

### Calculate Cart Price and Discounts

Add an endpoint to calculate the total price of the cart and apply discounts.
The total price is the sum of prices of all items in the cart.

Discount rules:
  - If the total price > 5000 → apply a 10% discount.
  - If the total number of items > 3 → apply a 5% discount.
  - If both conditions are met, apply the larger discount (10%).

Example request:

```sh
GET http://localhost:3000/carts/1/price
```

Example response:

```json
{
  "cart_id": 1,
  "total_price": 6200.00,
  "discount_percent": 10,
  "final_price": 5580.00
}
```

