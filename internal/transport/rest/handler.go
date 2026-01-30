package rest

import (
	"cart-api/internal/model"
	"cart-api/internal/repository/Cart"
	"cart-api/internal/services"
	"cart-api/internal/transport/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type CartProvider interface {
	CreateCart(context.Context) (*model.Cart, error)
	CreateItem(context.Context, model.CartItem) (int, error)
	DeleteItem(context.Context, model.CartItem) error
	GetCart(context.Context, int) (*model.Cart, error)
	GetPrice(context.Context, int) (*model.Price, error)
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
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	cartID := r.PathValue("cart_id")
	cartItem := r.PathValue("item_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Warn("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("invalid cart ID; '%s'must be an integer", cartID), http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(cartItem)
	if err != nil {
		h.logger.Warn("failed to parse cart item", zap.Error(err), zap.String("input", cartItem))
		http.Error(w, fmt.Sprintf("invalid item ID; '%s' must be an integer", cartItem), http.StatusBadRequest)
		return
	}
	item := model.CartItem{Id: itemID, CartId: id}
	err = h.service.DeleteItem(ctx, item)
	if err != nil {
		if errors.Is(err, Cart.ErrNotFound) {
			h.logger.Warn("attempt to delete non-existent item",
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
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	cart, err := h.service.CreateCart(ctx)
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
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Warn("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("invalid cart ID; '%s' must be an integer", cartID), http.StatusBadRequest)
		return
	}
	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	itemModel := model.CartItem{
		CartId:  id,
		Product: req.Product,
		Price:   req.Price,
	}
	newID, err := h.service.CreateItem(ctx, itemModel)
	if err != nil {

		if errors.Is(err, services.ErrInvalidProduct) ||
			errors.Is(err, services.ErrInvalidPrice) ||
			errors.Is(err, services.ErrReachCartLimit) {
			h.logger.Warn("business rule violation",
				zap.Error(err),
				zap.Int("cart_id", id),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, services.ErrCartNotFound) {
			h.logger.Warn("cart not found",
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
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Warn("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("invalid cart ID; '%s' must be an integer", cartID), http.StatusBadRequest)
		return
	}
	carts, err := h.service.GetCart(ctx, id)
	if err != nil {
		var notFoundErr *Cart.ErrCartNotFound
		if errors.As(err, &notFoundErr) {
			h.logger.Warn("cart not found",
				zap.Int("cart_id", id),
			)
			http.Error(w, notFoundErr.Error(), http.StatusNotFound)
			return
		}
		var notFoundItemErr *Cart.ErrCartItemNotFound
		if errors.As(err, &notFoundItemErr) {
			h.logger.Warn("cart items not found",
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
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	cartID := r.PathValue("cart_id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		h.logger.Warn("failed to parse cart id", zap.Error(err), zap.String("input", cartID))
		http.Error(w, fmt.Sprintf("Invalid cart_id '%s': must be an integer", cartID), http.StatusBadRequest)
		return
	}
	price, err := h.service.GetPrice(ctx, id)
	if err != nil {
		var notFoundErr *Cart.ErrCartNotFound
		if errors.As(err, &notFoundErr) {
			h.logger.Warn("cart not found for price calculation", zap.Int("id", id))
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
