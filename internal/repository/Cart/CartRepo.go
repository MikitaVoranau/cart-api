package Cart

import (
	"cart-api/internal/model/CartItem"
	"cart-api/internal/model/Carts"
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

func (r *CartRepo) CreateCart() (*Carts.Carts, error) {
	cart := &Carts.Carts{
		Items: []CartItem.CartItem{},
	}
	err := r.DB.QueryRow("INSERT INTO carts DEFAULT VALUES RETURNING id").Scan(&cart.ID)
	if err != nil {
		return nil, fmt.Errorf("error inserting carts: %w", err)
	}
	return cart, nil
}

func (r *CartRepo) CreateItem(item CartItem.CartItem) (int, error) {
	err := r.DB.QueryRow("SELECT add_item_to_cart ($1, $2, $3)", item.CartId, item.Product, item.Price).Scan(&item.Id)
	if err != nil {
		return 0, err
	}
	return item.Id, nil
}

func (r *CartRepo) DeleteItem(item CartItem.CartItem) error {
	res, err := r.DB.Exec("DELETE FROM cart_item WHERE id = $1 AND cart_id = $2", item.Id, item.CartId)
	if err != nil {
		return fmt.Errorf("could not delete item: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CartRepo) GetCart(id int) (*Carts.Carts, error) {
	carts := &Carts.Carts{}
	err := r.DB.QueryRow("SELECT id FROM carts WHERE id = $1", id).Scan(&carts.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &ErrCartNotFound{id}
		}
		return nil, fmt.Errorf("GetCart: query cart error: %w", err)
	}
	rows, err := r.DB.Query("SELECT * FROM cart_item WHERE cart_id = $1", carts.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &ErrCartItemNotFound{id, carts.ID}
		}
		return nil, fmt.Errorf("GetCart: query cart item error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item CartItem.CartItem
		if err = rows.Scan(&item.Id, &item.CartId, &item.Product, &item.Price); err != nil {
			return nil, fmt.Errorf("GetCart: scan item error: %w", err)
		}
		carts.Items = append(carts.Items, item)
	}
	return carts, nil
}
