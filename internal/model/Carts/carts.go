package Carts

import "cart-api/internal/model/CartItem"

type Carts struct {
	ID    int                 `json:"id"`
	Items []CartItem.CartItem `json:"items"`
}
