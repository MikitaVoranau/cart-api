package rest

import (
	"cart-api/internal/model"
	"cart-api/internal/repository/Cart"
	"cart-api/internal/transport/dto"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type CartProvider interface {
	CreateCart() (*model.Cart, error)
	CreateItem(item model.CartItem) (int, error)
	DeleteItem(item model.CartItem) error
	GetCart(id int) (*model.Cart, error)
	GetPrice(id int) (*model.Price, error)
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
		h.logger.Info("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("invalid cart ID; '%s'must be an integer", cartID), http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(cartItem)
	if err != nil {
		h.logger.Info("failed to parse cart item", zap.Error(err), zap.String("input", cartItem))
		http.Error(w, fmt.Sprintf("invalid item ID; '%s' must be an integer", cartItem), http.StatusBadRequest)
		return
	}
	item := model.CartItem{Id: itemID, CartId: id}
	err = h.service.DeleteItem(item)
	if err != nil {
		if errors.Is(err, Cart.ErrNotFound) {
			h.logger.Info("attempt to delete non-existent item",
				zap.Int("cart_id", id),
				zap.Int("item_id", itemID),
			)
			http.Error(w, fmt.Sprintf("Item %d not found in cart %d", itemID, id), http.StatusNotFound)
			return
		}
		h.logger.Error("failed to delete item",
			zap.Error(err),
			zap.Int("cart_id", id),
			zap.Int("item_id", itemID),
		)
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}
	h.logger.Info("item deleted successfully",
		zap.Int("cart_id", id),
		zap.Int("item_id", itemID),
	)
	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) PostCart(w http.ResponseWriter, r *http.Request) {
	cart, err := h.service.CreateCart()
	if err != nil {
		h.logger.Error("error creating cart", zap.Error(err))
		http.Error(w, "Failed to create new cart", http.StatusInternalServerError)
		return
	}
	resp := dto.CartResponse{
		ID:    cart.ID,
		Items: []dto.ItemResponse{},
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("error encoding cart", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CartHandler) PostItem(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Info("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("invalid cart ID; '%s' must be an integer", cartID), http.StatusBadRequest)
		return
	}
	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Info("invalid request body", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	itemModel := model.CartItem{
		CartId:  id,
		Product: req.Product,
		Price:   req.Price,
	}
	newID, err := h.service.CreateItem(itemModel)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "cart cannot consist more than 5 distinct products") {
			h.logger.Info("business rule violation",
				zap.Error(err),
				zap.Int("cart_id", id),
			)
			http.Error(w, "Cart limit reached: max 5 distinct products", http.StatusBadRequest)
			return
		}

		if strings.Contains(errMsg, "product name cannot be blank") ||
			strings.Contains(errMsg, "incorrect price") {
			h.logger.Info("validation error",
				zap.Error(err),
				zap.Int("cart_id", id),
				zap.String("product", req.Product),
				zap.Float64("price", req.Price),
			)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		if strings.Contains(errMsg, "does not exist") {
			h.logger.Info("cart not found",
				zap.Int("cart_id", id),
			)
			http.Error(w, fmt.Sprintf("Cart with id %d not found", id), http.StatusNotFound)
			return
		}

		h.logger.Error("failed to create item",
			zap.Error(err),
			zap.Int("cart_id", id),
			zap.String("product", req.Product),
			zap.Float64("price", req.Price),
		)
		http.Error(w, "Internal server error processing item creation", http.StatusInternalServerError)
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CartHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Info("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("invalid cart ID; '%s' must be an integer", cartID), http.StatusBadRequest)
		return
	}
	carts, err := h.service.GetCart(id)
	if err != nil {
		var notFoundErr *Cart.ErrCartNotFound
		if errors.As(err, &notFoundErr) {
			h.logger.Info("cart not found",
				zap.Int("cart_id", id),
			)
			http.Error(w, notFoundErr.Error(), http.StatusNotFound)
			return
		}
		var notFoundItemErr *Cart.ErrCartItemNotFound
		if errors.As(err, &notFoundItemErr) {
			h.logger.Info("cart items not found",
				zap.Int("cart_id", id),
			)
			http.Error(w, notFoundItemErr.Error(), http.StatusNotFound)
			return
		}
		h.logger.Error("error getting carts",
			zap.Error(err),
			zap.Int("cart_id", id),
		)
		http.Error(w, "Failed to receive cart details", http.StatusInternalServerError)
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CartHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Info("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("Invalid cart_id '%s': must be an integer", cartID), http.StatusBadRequest)
		return
	}
	price, err := h.service.GetPrice(id)
	if err != nil {
		var notFoundErr *Cart.ErrCartNotFound
		if errors.As(err, &notFoundErr) {
			h.logger.Info("cart not found for price calculation", zap.Int("id", id))
			http.Error(w, fmt.Sprintf("Cart with id %d not found", id), http.StatusNotFound)
			return
		}
		h.logger.Error("error calculating price",
			zap.Error(err),
			zap.Int("cart_id", id),
		)
		http.Error(w, "Failed to calculate cart price", http.StatusInternalServerError)
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
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
