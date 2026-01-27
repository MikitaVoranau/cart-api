package services

import (
	"cart-api/internal/model"
	"fmt"
	"math"
)

type CartRepository interface {
	GetCart(id int) (*model.Cart, error)
	CreateCart() (*model.Cart, error)
	CreateItem(model.CartItem) (int, error)
	DeleteItem(model.CartItem) error
	CartExists(cartID int) (bool, error)
	ItemExists(itemID int) (bool, error)
}

type CartService struct {
	CartRepo CartRepository
}

func NewCartService(cartRepo CartRepository) *CartService {
	return &CartService{
		cartRepo,
	}
}

func (s *CartService) CreateCart() (*model.Cart, error) {
	return s.CartRepo.CreateCart()
}

func (s *CartService) CreateItem(item model.CartItem) (int, error) {
	exists, err := s.CartRepo.CartExists(item.CartId)
	if err != nil {
		return 0, fmt.Errorf("failed to check cart existence: %w", err)
	}
	if !exists {
		return 0, ErrCartNotFound
	}
	return s.CartRepo.CreateItem(item)
}

func (s *CartService) DeleteItem(item model.CartItem) error {
	exists, err := s.CartRepo.ItemExists(item.Id)
	if err != nil {
		return fmt.Errorf("failed to check item existence: %w", err)
	}
	if !exists {
		return ErrItemNotFound
	}
	return s.CartRepo.DeleteItem(item)
}

func (s *CartService) GetCart(id int) (*model.Cart, error) {
	cart, err := s.CartRepo.GetCart(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil
}

func (s *CartService) GetPrice(id int) (*model.Price, error) {
	carts, err := s.CartRepo.GetCart(id)
	if err != nil {
		return nil, fmt.Errorf("getting cart for price failed: %w", err)
	}
	price := &model.Price{}
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
