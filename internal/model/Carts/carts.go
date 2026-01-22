package Carts

import "cart-api/internal/model/CartItem"

type Carts struct {
	ID    int
	Items []CartItem.CartItem
}
