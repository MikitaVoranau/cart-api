package CartRepo

import (
	"cart-api/internal/model/CartItem"
	"cart-api/internal/model/Carts"
	"github.com/jmoiron/sqlx"
	"log"
)

type CartRepo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *CartRepo {
	return &CartRepo{db}
}

// Delete SELECT INSERT
func (r *CartRepo) CreateCart() Carts.Carts {
	result, err := r.DB.Exec("INSERT INTO carts DEFAULT VALUES RETURNING *")
	if err != nil {
		log.Fatal(err)
	}
	var cart Carts.Carts
	lastID, _ := result.LastInsertId()
	cart.ID = int(lastID)
	return cart
}

func (r *CartRepo) CreateItem(item *CartItem.CartItem) {
	err := r.DB.QueryRow("SELECT add_item_to_cart ($1, $2, $3)", item.CartId, item.Product, item.Price).Scan(&item.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *CartRepo) DeleteItem(item CartItem.CartItem) {
	_, err := r.DB.Exec("DELETE FROM cart_item WHERE id = $1 AND cart_id = $2", item.Id, item.CartId)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *CartRepo) GetCart(id int) *Carts.Carts {
	carts := &Carts.Carts{}
	err := r.DB.QueryRow("SELECT * FROM carts WHERE id = $1", id).Scan(&carts.ID)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := r.DB.Query("SELECT * FROM cart_item WHERE cart_id = $1", carts.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var item CartItem.CartItem
		if err = rows.Scan(&item.Id, &item.CartId, &item.Product, &item.Price); err != nil {
			log.Fatal(err)
		}
		carts.Items = append(carts.Items, item)
	}
	return carts
}
