package Cart

import (
	"cart-api/internal/model"
	"cart-api/internal/repository/dao"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CartRepo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *CartRepo {
	return &CartRepo{db}
}

func (r *CartRepo) CreateCart(ctx context.Context) (*model.Cart, error) {
	var cartDb dao.CartDb
	err := r.DB.QueryRowxContext(ctx, "INSERT INTO carts DEFAULT VALUES RETURNING id").Scan(&cartDb.ID)
	if err != nil {
		return nil, fmt.Errorf("error inserting carts: %w", err)
	}
	return cartDb.ToDomain(), nil
}

func (r *CartRepo) CreateItem(ctx context.Context, item model.CartItem) (int, error) {
	itemDb := dao.NewCartItemDb(item)
	var newID int
	err := r.DB.QueryRowxContext(ctx, "INSERT INTO cart_item (cart_id, product, price) VALUES ($1, $2, $3) RETURNING id", itemDb.CartID, itemDb.Product, itemDb.Price).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *CartRepo) DeleteItem(ctx context.Context, item model.CartItem) error {
	res, err := r.DB.ExecContext(ctx, "DELETE FROM cart_item WHERE id = $1 AND cart_id = $2", item.Id, item.CartId)
	if err != nil {
		return fmt.Errorf("could not delete item: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CartRepo) GetCart(ctx context.Context, id int) (*model.Cart, error) {
	var cartDb dao.CartDb
	err := r.DB.QueryRowxContext(ctx, "SELECT id FROM carts WHERE id = $1", id).Scan(&cartDb.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &ErrCartNotFound{id}
		}
		return nil, fmt.Errorf("GetCart: query cart error: %w", err)
	}
	cart := cartDb.ToDomain()
	rows, err := r.DB.QueryxContext(ctx, "SELECT id, cart_id, product, price FROM cart_item WHERE cart_id = $1", cart.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &ErrCartItemNotFound{id, cart.ID}
		}
		return nil, fmt.Errorf("GetCart: query cart item error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var itemDb dao.CartItemDb
		if err = rows.Scan(&itemDb.ID, &itemDb.CartID, &itemDb.Product, &itemDb.Price); err != nil {
			return nil, fmt.Errorf("GetCart: scan item error: %w", err)
		}
		cart.Items = append(cart.Items, itemDb.ToDomain())
	}
	return cart, nil
}

func (r *CartRepo) CartExists(ctx context.Context, cartID int) (bool, error) {
	var exists bool
	err := r.DB.QueryRowxContext(ctx, "SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)", cartID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("CartExists: cart isn't found: %w", err)
	}
	return exists, nil
}

func (r *CartRepo) ItemExists(ctx context.Context, itemID int) (bool, error) {
	var exists bool
	err := r.DB.QueryRowxContext(ctx, "SELECT EXISTS(SELECT 1 FROM cart_item WHERE id = $1)", itemID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ItemExists: item isn't exist	: %w", err)
	}
	return exists, nil
}
