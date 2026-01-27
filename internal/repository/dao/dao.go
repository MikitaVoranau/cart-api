package dao

import "cart-api/internal/model"

type CartDb struct {
	ID int `db:"id"`
}

type CartItemDb struct {
	ID      int     `db:"id"`
	CartID  int     `db:"cart_id"`
	Product string  `db:"product"`
	Price   float64 `db:"price"`
}

func (dbItem *CartItemDb) ToDomain() model.CartItem {
	return model.CartItem{
		Id:      dbItem.ID,
		CartId:  dbItem.CartID,
		Product: dbItem.Product,
		Price:   dbItem.Price,
	}
}

func NewCartItemDb(item model.CartItem) CartItemDb {
	return CartItemDb{
		ID:      item.Id,
		CartID:  item.CartId,
		Product: item.Product,
		Price:   item.Price,
	}
}

func (dbCart *CartDb) ToDomain() *model.Cart {
	return &model.Cart{
		ID:    dbCart.ID,
		Items: []model.CartItem{},
	}
}
