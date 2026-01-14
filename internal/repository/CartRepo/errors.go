package CartRepo

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("item not found")
)

type ErrCartNotFound struct {
	ID int
}

func (e *ErrCartNotFound) Error() string {
	return fmt.Sprintf("cart with id %d not found", e.ID)
}

type ErrCartItemNotFound struct {
	ID     int
	cartID int
}

func (e *ErrCartItemNotFound) Error() string {
	return fmt.Sprintf("cart item with id %d not found with id %d", e.cartID, e.ID)
}
