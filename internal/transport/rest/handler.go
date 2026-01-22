package rest

import (
	"cart-api/internal/model/CartItem"
	"cart-api/internal/model/Carts"
	"cart-api/internal/model/Price"
	"cart-api/internal/repository/Cart"
	"cart-api/internal/transport/dto"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type CartProvider interface {
	CreateCart() (*Carts.Carts, error)
	CreateItem(item CartItem.CartItem) (int, error)
	DeleteItem(item CartItem.CartItem) error
	GetCart(id int) (*Carts.Carts, error)
	GetPrice(id int) (*Price.Price, error)
}

type CartHandler struct {
	service CartProvider
	logger  *zap.Logger
}

func NewCartHandler(service CartProvider, l *zap.Logger) *CartHandler {
	return &CartHandler{
		service,
		l,
	}
}

func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	cartItem := r.PathValue("item_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Info("error converting cartID to int", zap.Error(err))
		http.Error(w, "incorrect ID of cart", http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(cartItem)
	if err != nil {
		h.logger.Info("error converting cartID item to int", zap.Error(err))
		http.Error(w, "incorrect ID of cart item", http.StatusBadRequest)
		return
	}
	item := CartItem.CartItem{Id: itemID, CartId: id}
	err = h.service.DeleteItem(item)
	if err != nil {
		if errors.Is(err, Cart.ErrNotFound) {
			h.logger.Info("client tried to delete non-exist item", zap.Error(err), zap.Int("cart id", id), zap.Int("cart item id", itemID))
			http.Error(w, "cannot delete item", http.StatusNotFound)
			return
		}
		h.logger.Error("error deleting item", zap.Error(err), zap.Int("cart id", id), zap.Int("cart item id", itemID))
		http.Error(w, "cannot delete item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) PostCart(w http.ResponseWriter, r *http.Request) {
	cart, err := h.service.CreateCart()
	if err != nil {
		h.logger.Error("error creating cart", zap.Error(err))
		http.Error(w, "error creating cart", http.StatusInternalServerError)
		return
	}
	resp := dto.CartResponse{
		ID:    cart.ID,
		Items: []dto.ItemResponse{},
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("error encoding cart", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *CartHandler) PostItem(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Info("error converting cartID to int", zap.Error(err))
		http.Error(w, "incorrect ID of cart", http.StatusBadRequest)
		return
	}
	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Info("invalid request body", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	itemModel := CartItem.CartItem{
		CartId:  id,
		Product: req.Product,
		Price:   req.Price,
	}
	newID, err := h.service.CreateItem(itemModel)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "cart cannot consist more than 5 distinct products") {
			h.logger.Info("business rule violation", zap.Error(err))
			http.Error(w, "Cart limit reached: max 5 distinct products", http.StatusBadRequest)
			return
		}

		if strings.Contains(errMsg, "product name cannot be blank") ||
			strings.Contains(errMsg, "incorrect price") {
			h.logger.Info("validation error", zap.Error(err))
			http.Error(w, errMsg, http.StatusBadRequest) // Можно отдать текст ошибки клиенту
			return
		}

		if strings.Contains(errMsg, "does not exist") {
			h.logger.Info("cart not found", zap.Int("cart_id", id))
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}

		h.logger.Error("failed to create item", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	resp := dto.ItemResponse{
		ID:      newID,
		CartID:  id,
		Product: itemModel.Product,
		Price:   itemModel.Price,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("error encoding cartItem", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *CartHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Info("error converting cartID to int", zap.Error(err))
		http.Error(w, "incorrect ID of cart", http.StatusBadRequest)
		return
	}
	carts, err := h.service.GetCart(id)
	if err != nil {
		var notFoundErr *Cart.ErrCartNotFound
		if errors.As(err, &notFoundErr) {
			h.logger.Info("cart not found", zap.Int("id", notFoundErr.ID))
			http.Error(w, notFoundErr.Error(), http.StatusNotFound)
			return
		}
		var notFoundItemErr *Cart.ErrCartItemNotFound
		if errors.As(err, &notFoundItemErr) {
			h.logger.Info("cart items not found", zap.Int("id", notFoundErr.ID), zap.Int("cart_id", id))
			http.Error(w, notFoundItemErr.Error(), http.StatusNotFound)
			return
		}
		h.logger.Error("error getting carts", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	itemsDTO := make([]dto.ItemResponse, 0, len(carts.Items))

	for _, item := range carts.Items {
		itemsDTO = append(itemsDTO, dto.ItemResponse{
			ID:      item.Id,
			CartID:  item.CartId,
			Product: item.Product,
			Price:   item.Price,
		})
	}

	resp := dto.CartResponse{
		ID:    carts.ID,
		Items: itemsDTO,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("error encoding carts", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *CartHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Info("error converting cartID to int", zap.Error(err))
		http.Error(w, "incorrect ID of cart", http.StatusBadRequest)
		return
	}
	price, err := h.service.GetPrice(id)
	var notFoundErr *Cart.ErrCartNotFound
	if errors.As(err, &notFoundErr) {
		h.logger.Info("cart not found for price calculation", zap.Int("id", id))
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}
	resp := dto.PriceResponse{
		CartID:          price.CartId,
		TotalPrice:      price.TotalPrice,
		DiscountPercent: price.DiscountPercent,
		FinalPrice:      price.FinalPrice,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("error encoding price", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
