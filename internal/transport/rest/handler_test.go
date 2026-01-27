package rest

import (
	"bytes"
	"cart-api/internal/model"
	"cart-api/internal/repository/Cart"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateCart() (*model.Cart, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockService) CreateItem(item model.CartItem) (int, error) {
	args := m.Called(item)
	return args.Int(0), args.Error(1)
}

func (m *MockService) DeleteItem(item model.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockService) GetCart(id int) (*model.Cart, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockService) GetPrice(id int) (*model.Price, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Price), args.Error(1)
}

func TestCartHandler_PostCart(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockSvc := new(MockService)
	handler := NewCartHandler(mockSvc, logger)

	t.Run("Success", func(t *testing.T) {
		expectedCart := &model.Cart{ID: 100}
		mockSvc.On("CreateCart").Return(expectedCart, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/carts", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("POST /carts", handler.PostCart)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":100`)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc.On("CreateCart").Return(nil, errors.New("db fail")).Once()

		req := httptest.NewRequest(http.MethodPost, "/carts", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("POST /carts", handler.PostCart)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestCartHandler_PostItem(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockSvc := new(MockService)
	handler := NewCartHandler(mockSvc, logger)

	validItem := model.CartItem{Product: "Apple", Price: 50}
	bodyJSON, _ := json.Marshal(validItem)

	tests := []struct {
		name           string
		cartID         string
		body           []byte
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success",
			cartID: "1",
			body:   bodyJSON,
			setupMock: func() {
				mockSvc.On("CreateItem", mock.MatchedBy(func(i model.CartItem) bool {
					return i.Product == "Apple" && i.CartId == 1
				})).Return(555, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":555`,
		},
		{
			name:   "Invalid Cart ID",
			cartID: "abc",
			body:   bodyJSON,
			setupMock: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Business Limit Error",
			cartID: "1",
			body:   bodyJSON,
			setupMock: func() {
				mockSvc.On("CreateItem", mock.Anything).Return(0, errors.New("cart cannot consist more than 5 distinct products"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "limit reached",
		},
		{
			name:   "Not Found Error",
			cartID: "99",
			body:   bodyJSON,
			setupMock: func() {
				mockSvc.On("CreateItem", mock.Anything).Return(0, errors.New("something does not exist"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Internal Error",
			cartID: "1",
			body:   bodyJSON,
			setupMock: func() {
				mockSvc.On("CreateItem", mock.Anything).Return(0, errors.New("unknown db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/carts/"+tt.cartID+"/items", bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("POST /carts/{cart_id}/items", handler.PostItem)
			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
			mockSvc.AssertExpectations(t)
			mockSvc.ExpectedCalls = nil
		})
	}
}

func TestCartHandler_GetItems(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockSvc := new(MockService)
	handler := NewCartHandler(mockSvc, logger)

	t.Run("Success", func(t *testing.T) {
		cart := &model.Cart{ID: 1, Items: []model.CartItem{{Product: "Banana"}}}
		mockSvc.On("GetCart", 1).Return(cart, nil)

		req := httptest.NewRequest(http.MethodGet, "/carts/1", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /carts/{cart_id}", handler.GetItems)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Banana")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Not Found (Custom Error Type)", func(t *testing.T) {
		notFoundErr := &Cart.ErrCartNotFound{ID: 999}
		mockSvc.On("GetCart", 999).Return(nil, notFoundErr)

		req := httptest.NewRequest(http.MethodGet, "/carts/999", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /carts/{cart_id}", handler.GetItems)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestCartHandler_DeleteItem(t *testing.T) {
	logger := zaptest.NewLogger(t)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockService)
		handler := NewCartHandler(mockSvc, logger)

		mockSvc.On("DeleteItem", mock.MatchedBy(func(i model.CartItem) bool {
			return i.CartId == 1 && i.Id == 5
		})).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/carts/1/items/5", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /carts/{cart_id}/items/{item_id}", handler.DeleteItem)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Not Found (Var Error)", func(t *testing.T) {
		mockSvc := new(MockService)
		handler := NewCartHandler(mockSvc, logger)

		mockSvc.On("DeleteItem", mock.Anything).Return(Cart.ErrNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/carts/1/items/5", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /carts/{cart_id}/items/{item_id}", handler.DeleteItem)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestCartHandler_GetPrice(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockSvc := new(MockService)
	handler := NewCartHandler(mockSvc, logger)

	t.Run("Success", func(t *testing.T) {
		price := &model.Price{FinalPrice: 1000}
		mockSvc.On("GetPrice", 1).Return(price, nil)

		req := httptest.NewRequest(http.MethodGet, "/carts/1/price", nil)
		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /carts/{cart_id}/price", handler.GetPrice)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "1000")
		mockSvc.AssertExpectations(t)
	})
}
