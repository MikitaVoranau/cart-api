package services

import (
	"errors"
)

var (
	ErrCartNotFound   = errors.New("cart not found")
	ErrItemNotFound   = errors.New("item not found")
	ErrInvalidProduct = errors.New("product name cannot be blank")
	ErrInvalidPrice   = errors.New("incorrect price information")
	ErrReachCartLimit = errors.New("cart limit reached: max 5 distinct products")
)
