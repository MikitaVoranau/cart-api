package rest

import (
	"cart-api/internal/model/CartItem"
	"cart-api/internal/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{
		service,
	}
}

func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	cartItem := r.PathValue("item_id")
	id, _ := strconv.Atoi(cartID)
	itemID, _ := strconv.Atoi(cartItem)
	fmt.Println(itemID)
	item := CartItem.CartItem{Id: id, CartId: itemID}
	h.service.CartRepo.DeleteItem(item)
	err := json.NewEncoder(w).Encode(item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *CartHandler) PostCart(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	cart := h.service.CartRepo.CreateCart()
	err := json.NewEncoder(w).Encode(cart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *CartHandler) PostItem(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, _ := strconv.Atoi(cartID)
	var cartItem *CartItem.CartItem
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, cartItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cartItem.CartId = id
	h.service.CartRepo.CreateItem(cartItem)
}

func (h *CartHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, _ := strconv.Atoi(cartID)
	carts := h.service.CartRepo.GetCart(id)
	err := json.NewEncoder(w).Encode(carts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *CartHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, _ := strconv.Atoi(cartID)
	price := h.service.GetPrice(id)
	err := json.NewEncoder(w).Encode(price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
