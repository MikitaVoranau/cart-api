package model

type CartItem struct {
	Id      int
	CartId  int
	Product string
	Price   float64
}
type Cart struct {
	ID    int
	Items []CartItem
}

type Price struct {
	CartId          int
	TotalPrice      float64
	DiscountPercent int
	FinalPrice      float64
}
