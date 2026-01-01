package CartItem

type CartItem struct {
	Id      int     `json:"id"`
	CartId  int     `json:"cart_id"`
	Product string  `json:"product"`
	Price   float64 `json:"price"`
}
