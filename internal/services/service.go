package services

import (
	"cart-api/internal/model/CartItem"
	"cart-api/internal/model/Carts"
	"cart-api/internal/model/Price"
	"fmt"
	"math"
)

type CartRepository interface {
	GetCart(id int) (*Carts.Carts, error)
	CreateCart() (*Carts.Carts, error)
	CreateItem(item CartItem.CartItem) (int, error)
	DeleteItem(CartItem.CartItem) error
}

type CartService struct {
	CartRepo CartRepository
}

func NewCartService(cartRepo CartRepository) *CartService {
	return &CartService{
		cartRepo,
	}
}

func (s *CartService) CreateCart() (*Carts.Carts, error) {
	return s.CartRepo.CreateCart()
}

func (s *CartService) CreateItem(item CartItem.CartItem) (int, error) {
	return s.CartRepo.CreateItem(item)
}

func (s *CartService) DeleteItem(item CartItem.CartItem) error {
	return s.CartRepo.DeleteItem(item)
}

func (s *CartService) GetCart(id int) (*Carts.Carts, error) {
	return s.CartRepo.GetCart(id)
}

func (cartService *CartService) GetPrice(id int) (*Price.Price, error) {
	carts, err := cartService.CartRepo.GetCart(id)
	if err != nil {
		return nil, fmt.Errorf("getting cart for price failed: %w", err)
	}
	price := &Price.Price{}
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
	price.FinalPrice = math.Trunc((totalPrice-totalPrice*(float64(price.DiscountPercent)/100))*100) / 100
	price.CartId = carts.ID
	return price, nil
}
