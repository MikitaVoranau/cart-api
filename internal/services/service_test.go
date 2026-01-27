package services

import (
	"cart-api/internal/model"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCartRepo struct {
	mock.Mock
}

func (m *MockCartRepo) GetCart(id int) (*model.Cart, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockCartRepo) CreateCart() (*model.Cart, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockCartRepo) CreateItem(item model.CartItem) (int, error) {
	args := m.Called(item)
	return args.Int(0), args.Error(1)
}

func (m *MockCartRepo) DeleteItem(item model.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockCartRepo) CartExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockCartRepo) ItemExists(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestCreateCart(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		expectedCart := &model.Cart{ID: 1, Items: []model.CartItem{}}

		mockRepo.On("CreateCart").Return(expectedCart, nil)

		service := NewCartService(mockRepo)
		result, err := service.CreateCart()

		assert.NoError(t, err)
		assert.Equal(t, expectedCart, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		mockRepo.On("CreateCart").Return(nil, errors.New("db fail"))

		service := NewCartService(mockRepo)
		result, err := service.CreateCart()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateItem(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		item := model.CartItem{CartId: 1, Product: "Apple", Price: 100}
		expectedID := 123

		mockRepo.On("CartExists", item.CartId).Return(true, nil)
		mockRepo.On("CreateItem", item).Return(expectedID, nil)

		service := NewCartService(mockRepo)
		id, err := service.CreateItem(item)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Cart Not Found", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		item := model.CartItem{CartId: 999, Product: "Apple"}

		mockRepo.On("CartExists", item.CartId).Return(false, nil)

		service := NewCartService(mockRepo)
		id, err := service.CreateItem(item)

		assert.Error(t, err)
		assert.Equal(t, ErrCartNotFound, err)
		assert.Zero(t, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repo Error on Create", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		item := model.CartItem{CartId: 1, Product: "Apple"}

		mockRepo.On("CartExists", item.CartId).Return(true, nil)
		mockRepo.On("CreateItem", item).Return(0, errors.New("insert failed"))

		service := NewCartService(mockRepo)
		id, err := service.CreateItem(item)

		assert.Error(t, err)
		assert.Zero(t, id)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteItem(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		item := model.CartItem{Id: 10, CartId: 5}

		mockRepo.On("ItemExists", item.Id).Return(true, nil)
		mockRepo.On("DeleteItem", item).Return(nil)

		service := NewCartService(mockRepo)
		err := service.DeleteItem(item)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		item := model.CartItem{Id: 10, CartId: 5}

		mockRepo.On("ItemExists", item.Id).Return(false, nil)

		service := NewCartService(mockRepo)
		err := service.DeleteItem(item)

		assert.Error(t, err)
		assert.Equal(t, ErrItemNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetCart(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		expectedCart := &model.Cart{ID: 55}

		mockRepo.On("GetCart", 55).Return(expectedCart, nil)

		service := NewCartService(mockRepo)
		cart, err := service.GetCart(55)

		assert.NoError(t, err)
		assert.Equal(t, expectedCart, cart)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(MockCartRepo)
		mockRepo.On("GetCart", 55).Return(nil, errors.New("db error"))

		service := NewCartService(mockRepo)
		cart, err := service.GetCart(55)

		assert.Error(t, err)
		assert.Nil(t, cart)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPrice(t *testing.T) {
	tests := []struct {
		name           string
		cartID         int
		mockReturnCart *model.Cart
		mockReturnErr  error
		expectedPrice  float64
		expectedDisc   int
		expectError    bool
	}{
		{
			name:   "No Discount",
			cartID: 1,
			mockReturnCart: &model.Cart{
				ID: 1,
				Items: []model.CartItem{
					{Price: 100}, {Price: 200},
				},
			},
			mockReturnErr: nil,
			expectedPrice: 300.0,
			expectedDisc:  0,
			expectError:   false,
		},
		{
			name:   "Quantity Discount 5%",
			cartID: 2,
			mockReturnCart: &model.Cart{
				ID: 2,
				Items: []model.CartItem{
					{Price: 100}, {Price: 100}, {Price: 100}, {Price: 100},
				},
			},
			mockReturnErr: nil,
			expectedPrice: 380.0,
			expectedDisc:  5,
			expectError:   false,
		},
		{
			name:   "Amount Discount 10%",
			cartID: 3,
			mockReturnCart: &model.Cart{
				ID: 3,
				Items: []model.CartItem{
					{Price: 6000},
				},
			},
			mockReturnErr: nil,
			expectedPrice: 5400.0,
			expectedDisc:  10,
			expectError:   false,
		},
		{
			name:           "Repo Error",
			cartID:         5,
			mockReturnCart: nil,
			mockReturnErr:  errors.New("db error"),
			expectedPrice:  0,
			expectedDisc:   0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCartRepo)

			mockRepo.On("GetCart", tt.cartID).Return(tt.mockReturnCart, tt.mockReturnErr)

			service := NewCartService(mockRepo)
			gotPrice, err := service.GetPrice(tt.cartID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotPrice)
				assert.Equal(t, tt.expectedPrice, gotPrice.FinalPrice)
				assert.Equal(t, tt.expectedDisc, gotPrice.DiscountPercent)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
