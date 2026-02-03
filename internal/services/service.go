package services

import (
	"cart-api/internal/model"
	"context"
	"fmt"
	"math"
	"strings"
)

type CartRepository interface {
	GetCart(context.Context, int) (*model.Cart, error)
	CreateCart(context.Context) (*model.Cart, error)
	CreateItem(context.Context, model.CartItem) (int, error)
	DeleteItem(context.Context, model.CartItem) error
	CartExists(context.Context, int) (bool, error)
	ItemExists(context.Context, int) (bool, error)
}

type CartService struct {
	CartRepo CartRepository
}

func NewCartService(cartRepo CartRepository) *CartService {
	return &CartService{
		cartRepo,
	}
}

func (s *CartService) CreateCart(ctx context.Context) (*model.Cart, error) {
	return s.CartRepo.CreateCart(ctx)
}

func (s *CartService) CreateItem(ctx context.Context, item model.CartItem) (int, error) {
	exists, err := s.CartRepo.CartExists(ctx, item.CartId)
	if err != nil {
		return 0, fmt.Errorf("failed to check cart existence: %w", err)
	}
	if !exists {
		return 0, ErrCartNotFound
	}
	if strings.TrimSpace(item.Product) == "" {
		return 0, ErrInvalidProduct
	}
	if item.Price < 0 {
		return 0, ErrInvalidPrice
	}
	cart, err := s.GetCart(ctx, item.CartId)
	if err != nil {
		return 0, fmt.Errorf("failed to get cart: %w", err)
	}
	uniqueProducts := make(map[string]struct{})
	productAlreadyExist := false
	for _, existingItem := range cart.Items {
		uniqueProducts[existingItem.Product] = struct{}{}
		if existingItem.Product == item.Product {
			productAlreadyExist = true
		}
	}
	if len(uniqueProducts) >= 5 && !productAlreadyExist {
		return 0, ErrReachCartLimit
	}
	return s.CartRepo.CreateItem(ctx, item)
}

func (s *CartService) DeleteItem(ctx context.Context, item model.CartItem) error {
	exists, err := s.CartRepo.ItemExists(ctx, item.Id)
	if err != nil {
		return fmt.Errorf("failed to check item existence: %w", err)
	}
	if !exists {
		return ErrItemNotFound
	}
	return s.CartRepo.DeleteItem(ctx, item)
}

func (s *CartService) GetCart(ctx context.Context, id int) (*model.Cart, error) {
	cart, err := s.CartRepo.GetCart(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil
}

func (s *CartService) GetPrice(ctx context.Context, id int) (*model.Price, error) {
	carts, err := s.CartRepo.GetCart(ctx, id)
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
