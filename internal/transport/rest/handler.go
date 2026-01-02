package rest

import (
	"cart-api/internal/model/CartItem"
	"cart-api/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{
		service,
	}
}

func (h *CartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	switch r.Method {
	case http.MethodGet:
		fmt.Println(params)
		switch len(params) {
		case 1:
			cartID, err := strconv.Atoi(params[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			carts := h.service.CartRepo.GetCart(cartID)
			err = json.NewEncoder(w).Encode(carts)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case 2:
			cartID, err := strconv.Atoi(params[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if params[1] != "price" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			price := h.service.GetPrice(cartID)
			err = json.NewEncoder(w).Encode(price)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	case http.MethodPost:
		switch len(params) {
		case 0:
			cart := h.service.CartRepo.CreateCart()
			err := json.NewEncoder(w).Encode(cart)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case 2:
			cartID, err := strconv.Atoi(params[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var cartItem *CartItem.CartItem
			err = json.Unmarshal([]byte(params[1]), cartItem)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			cartItem.CartId = cartID
			h.service.CartRepo.CreateItem(cartItem)
		}
	case http.MethodDelete:
		cartID, err := strconv.Atoi(params[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cartItem, err := strconv.Atoi(params[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		item := CartItem.CartItem{Id: cartItem, CartId: cartID}
		h.service.CartRepo.DeleteItem(item)
		err = json.NewEncoder(w).Encode(item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}
