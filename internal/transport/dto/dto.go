package dto

type AddItemRequest struct {
	Product string  `json:"product"`
	Price   float64 `json:"price"`
}

type ItemResponse struct {
	ID      int     `json:"id"`
	CartID  int     `json:"cart_id"`
	Product string  `json:"product"`
	Price   float64 `json:"price"`
}

type DeleteItemRequest struct {
	ID     int `json:"id"`
	CartID int `json:"cart_id"`
}

type CartResponse struct {
	ID    int            `json:"id"`
	Items []ItemResponse `json:"items"`
}

type PriceResponse struct {
	CartID          int     `json:"cart_id"`
	TotalPrice      float64 `json:"total_price"`
	DiscountPercent int     `json:"discount_percent"`
	FinalPrice      float64 `json:"final_price"`
}
