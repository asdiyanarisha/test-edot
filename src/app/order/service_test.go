package order

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"test-edot/constants"
	"test-edot/src/dto"
	"test-edot/src/models"
	"test-edot/src/repository/mocks"
	"testing"
)

type (
	TestProcessOrderData struct {
		name             string
		err              error
		message          string
		payload          dto.PayloadCreateOrder
		mockResponse     []models.StockLevelProduct
		expectTotalOrder int
		expectTotalPrice float64
		isErr            bool
	}
)

func TestProcessOrder(t *testing.T) {

	tableTests := []TestProcessOrderData{
		{
			name:    "test normal order",
			err:     nil,
			message: "please enter startDate endDate valid",
			payload: dto.PayloadCreateOrder{Items: []dto.PayloadCreateOrderItems{
				{
					ProductId: 1,
					Qty:       2,
				},
			}},
			mockResponse: []models.StockLevelProduct{
				{ProductId: 1, WarehouseId: 1, Stock: 3, ReservedStock: 0, Product: models.Product{Price: 10000}},
			},
			expectTotalOrder: 1,
			isErr:            false,
			expectTotalPrice: 20000,
		},
		{
			name:    "test normal order with double warehouse",
			err:     nil,
			message: "please enter startDate endDate valid",
			payload: dto.PayloadCreateOrder{Items: []dto.PayloadCreateOrderItems{
				{
					ProductId: 2,
					Qty:       4,
				},
			}},
			mockResponse: []models.StockLevelProduct{
				{ProductId: 2, WarehouseId: 1, Stock: 2, ReservedStock: 0, Product: models.Product{Price: 3000}},
				{ProductId: 2, WarehouseId: 2, Stock: 2, ReservedStock: 0, Product: models.Product{Price: 3000}},
			},
			expectTotalOrder: 2,
			isErr:            false,
			expectTotalPrice: 12000,
		},
		{
			name:    "test error doesnt have stock",
			err:     constants.NotEnoughStockProduct,
			message: "please enter startDate endDate valid",
			payload: dto.PayloadCreateOrder{Items: []dto.PayloadCreateOrderItems{
				{
					ProductId: 2,
					Qty:       4,
				},
			}},
			mockResponse: []models.StockLevelProduct{
				{ProductId: 2, WarehouseId: 1, Stock: 0, ReservedStock: 0, Product: models.Product{Price: 3000}},
			},
			expectTotalOrder: 0,
			isErr:            true,
			expectTotalPrice: 0,
		},
	}

	for _, test := range tableTests {
		t.Run(test.name, func(t *testing.T) {

			tx := gorm.DB{}
			mockRepo := new(mocks.StockLevelRepositoryInterface)

			mockRepo.On("FindTx", &tx, "updated_at asc", "product_id = ? and stock > 0", test.payload.Items[0].ProductId).Return(test.mockResponse, nil)

			mockRepo.On("UpdateOneTx", &tx, mock.AnythingOfType("*models.StockLevel"), mock.Anything, mock.Anything, 0).Return(nil)

			s := service{StockLevelRepository: mockRepo}

			order, totalPrice, err := s.ProcessOrder(&tx, test.payload)
			if test.isErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectTotalOrder, len(order))
				assert.Equal(t, test.expectTotalPrice, totalPrice)
			}
		})
	}
}
