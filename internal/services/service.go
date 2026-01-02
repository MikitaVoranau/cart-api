package services

import (
	"cart-api/internal/model/Price"
	"cart-api/internal/repository/CartRepo"
)

type CartService struct {
	CartRepo *CartRepo.CartRepo
}

func NewCartService(cartRepo *CartRepo.CartRepo) *CartService {
	return &CartService{
		cartRepo,
	}
}

func (cartService *CartService) GetPrice(id int) Price.Price {
	carts := cartService.CartRepo.GetCart(id)

	price := Price.Price{}
	var totalPrice float64
	var totalNumbers int

	for _, currPrice := range carts.Items {
		totalNumbers++
		totalPrice += currPrice.Price
	}

	price.TotalPrice = totalPrice
	if totalNumbers > 3 {
		price.DiscountPercent = 5
	}
	if totalPrice > 5000 {
		price.DiscountPercent = 10
	}
	price.FinalPrice = totalPrice * (float64(price.DiscountPercent) / 100)
	price.CartId = carts.ID
	return price
}
